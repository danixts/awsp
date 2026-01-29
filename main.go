package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/danixts/awsp/internal/aws"
	"github.com/danixts/awsp/internal/install"
	"github.com/danixts/awsp/internal/msg"
	"github.com/danixts/awsp/internal/selector"
	"github.com/danixts/awsp/internal/store"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	flagValidate bool
	flagFull     bool
)

var rootCmd = &cobra.Command{
	Use:   "awsp",
	Short: "AWS Profile Switcher",
	Long:  "CLI to switch between AWS profiles and regions. Run without args for interactive mode.",
	RunE:  runInteractive,
}

var switchCmd = &cobra.Command{
	Use:   "switch [profile] [region]",
	Short: "Switch profile (and optionally region)",
	Long:  "With no args: interactive selection. With profile: print export lines. With profile + region: set both.",
	Args:  cobra.MaximumNArgs(2),
	RunE:  runSwitch,
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	RunE:  runList,
}

var currentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show current profile",
	RunE:  runCurrent,
}

var favoriteCmd = &cobra.Command{
	Use:   "favorite",
	Short: "Manage favorite profiles",
}

var favoriteAddCmd = &cobra.Command{
	Use:   "add [profile]",
	Short: "Add profile to favorites",
	Args:  cobra.ExactArgs(1),
	RunE:  runFavoriteAdd,
}

var favoriteRemoveCmd = &cobra.Command{
	Use:   "remove [profile]",
	Short: "Remove profile from favorites",
	Args:  cobra.ExactArgs(1),
	RunE:  runFavoriteRemove,
}

var favoriteListCmd = &cobra.Command{
	Use:   "list",
	Short: "List favorite profiles",
	RunE:  runFavoriteList,
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Build, install binary and configure shell (zsh/bash/Windows)",
	Long:  "Builds the binary, copies it to a PATH directory, then asks for your shell to add awsenv and completion.",
	RunE:  runInstall,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&flagFull, "full", "f", false, "Export all AWS vars (including AWS_SDK_LOAD_CONFIG)")
	rootCmd.PersistentFlags().BoolVarP(&flagValidate, "validate", "v", false, "Verify credentials with aws sts get-caller-identity")

	rootCmd.AddCommand(switchCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(favoriteCmd)
	rootCmd.AddCommand(installCmd)

	favoriteCmd.AddCommand(favoriteAddCmd)
	favoriteCmd.AddCommand(favoriteRemoveCmd)
	favoriteCmd.AddCommand(favoriteListCmd)
}

func loadProfilesOrExit() []aws.Profile {
	profiles, err := aws.LoadProfiles()
	if err != nil {
		fmt.Printf(msg.ErrLoadProfiles+"\n", err)
		os.Exit(1)
	}
	if len(profiles) == 0 {
		fmt.Println(msg.ErrNoProfiles)
		os.Exit(1)
	}
	return profiles
}

func curProfileAndRegion() (profile, region string) {
	profile = aws.CurrentProfile()
	region = os.Getenv("AWS_REGION")
	if region == "" {
		region = os.Getenv("AWS_DEFAULT_REGION")
	}
	if last, err := store.LoadLast(); err == nil {
		if profile == aws.DefaultProfile && last.Profile != "" {
			profile = last.Profile
		}
		if region == "" && last.Region != "" {
			region = last.Region
		}
	}
	return profile, region
}

func runInteractive(cmd *cobra.Command, args []string) error {
	if !isatty.IsTerminal(os.Stdout.Fd()) {
		last, err := store.LoadLast()
		if err != nil || last.Profile == "" {
			fmt.Fprintln(os.Stderr, "awsp: output is being captured (no TTY). Run 'awsp' in your terminal once to select a profile.")
			os.Exit(1)
		}
		printExports(last.Profile, last.Region)
		return nil
	}
	profiles := loadProfilesOrExit()
	curProfile, curRegion := curProfileAndRegion()

	favs, _ := store.LoadFavorites()
	profiles = aws.OrderWithFavoritesFirst(profiles, favs)

	name, err := selector.RunProfileSelector(profiles, curProfile)
	if err != nil {
		fmt.Println(msg.Canceled)
		os.Exit(0)
	}

	p, _ := aws.FindProfile(profiles, name)
	region := p.Region
	if region == "" {
		regions := aws.RegionListForSelector("")
		region, err = selector.RunRegionSelector(regions, curRegion)
		if err != nil {
			fmt.Println(msg.Canceled)
			os.Exit(0)
		}
	}

	if flagValidate {
		if err := aws.ValidateProfile(name); err != nil {
			fmt.Printf(msg.ValidateFailed, err)
			os.Exit(1)
		}
	}

	_ = store.SaveLast(name, region)
	printExports(name, region)
	return nil
}

func runSwitch(cmd *cobra.Command, args []string) error {
	profiles := loadProfilesOrExit()
	_, curRegion := curProfileAndRegion()

	if len(args) == 0 {
		return runInteractive(cmd, args)
	}

	name := args[0]
	p, found := aws.FindProfile(profiles, name)
	if !found {
		fmt.Printf(msg.ErrProfileNotFound+"\n", name)
		os.Exit(1)
	}

	region := p.Region
	if len(args) > 1 {
		region = args[1]
	} else if region == "" {
		region = curRegion
	}

	if flagValidate {
		if err := aws.ValidateProfile(name); err != nil {
			fmt.Printf(msg.ValidateFailed, err)
			os.Exit(1)
		}
	}

	_ = store.SaveLast(name, region)
	printExports(name, region)
	return nil
}

func exportLines(profile, region string) string {
	var b strings.Builder
	b.WriteString(msg.ExportClearCreds)
	fmt.Fprintf(&b, msg.ExportProfile, profile)
	if region != "" {
		fmt.Fprintf(&b, msg.ExportRegion, region, region)
	}
	if flagFull {
		b.WriteString(msg.ExportSDKLoadConfig)
	}
	return b.String()
}

func printExports(profile, region string) {
	lines := exportLines(profile, region)
	fmt.Print(lines)
	fmt.Print(msg.HintApply)
	if path := os.Getenv("AWSP_EXPORT_FILE"); path != "" {
		_ = os.WriteFile(path, []byte(lines), 0600)
	}
}

func runList(cmd *cobra.Command, args []string) error {
	profiles := loadProfilesOrExit()
	favs, _ := store.LoadFavorites()
	profiles = aws.OrderWithFavoritesFirst(profiles, favs)
	for _, p := range profiles {
		star := ""
		if store.IsFavorite(p.Name, favs) {
			star = " ★"
		}
		if p.Region != "" {
			fmt.Printf(msg.ListItemWithReg, p.Name+star, p.Region)
		} else {
			fmt.Printf(msg.ListItem, p.Name+star)
		}
	}
	return nil
}

func runCurrent(cmd *cobra.Command, args []string) error {
	cur := aws.CurrentProfile()
	fmt.Printf(msg.CurrentProfile, cur)
	return nil
}

func runFavoriteAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	if err := store.AddFavorite(name); err != nil {
		return err
	}
	fmt.Printf(msg.FavoriteAdded, name)
	return nil
}

func runFavoriteRemove(cmd *cobra.Command, args []string) error {
	name := args[0]
	if err := store.RemoveFavorite(name); err != nil {
		return err
	}
	fmt.Printf(msg.FavoriteRemoved, name)
	return nil
}

func runFavoriteList(cmd *cobra.Command, args []string) error {
	favs, err := store.LoadFavorites()
	if err != nil {
		return err
	}
	if len(favs) == 0 {
		fmt.Print(msg.NoFavorites)
		return nil
	}
	fmt.Print(msg.FavoritesList)
	for _, name := range favs {
		fmt.Printf(msg.ListItem, name)
	}
	return nil
}

func runInstall(cmd *cobra.Command, args []string) error {
	return install.Run()
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
