ver las apis gateways

aws apigateway get-rest-apis \
  --query "items[].{id:id,name:name}" \
  --output table

masivo-integration-apigateway-dev-

aws apigateway get-resources \
  --rest-api-id x46efqn7vh \
  --query "items[].{id:id,path:path,methods:keys(resourceMethods)}" \
  --output table

aws apigateway get-resources \
  --rest-api-id x46efqn7vh \
  --query "items[?resourceMethods!=null].{path:path,methods:keys(resourceMethods)}"
  --output table

aws apigateway get-resources \
  --rest-api-id x46efqn7vh \
  --query "items[?resourceMethods!=null].{path:path,methods:join(',', keys(resourceMethods))}" \
  --output table

aws logs tail "/aws/lambda/masivo-integration-plugin-dev-PostContactsFunction-ROhl0dQytQeJ" \
  --since 1 \
  --follow \
  --region us-east-1