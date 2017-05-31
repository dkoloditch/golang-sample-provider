# golang-sample-provider

This is a minimal provider for Manifold written in Go. It's similar to [this one written in Sinatra](https://github.com/manifoldco/ruby-sinatra-sample-provider) and returns random numbers.

# Testing the app with Grafton

[Grafton](https://github.com/manifoldco/grafton) is a test framework used to verify Manifold provider implementations like this one.

Contact [support@manifold.co](mailto:support@manifold.co) for access to
Grafton and [check out the docs here](https://docs.manifold.co/#section/Getting-Started/Using-Grafton) for more info.


To use Grafton to verify this sample provider:

```bash
# Create a test master key for grafton to use when acting as Manifold
# this file is written as masterkey.json.
grafton generate

# Set the following environment variables to configure the test app.
# MASTER_KEY: the public_key portion of masterkey.json
# CONNECTOR_URL: the url that Grafton will listen on. It corresponds to
#   Grafton's --sso-port flag.
# CLIENT_ID & CLIENT_SECRET: fake OAuth 2.0 credentials. The format of these
#   are specific, so you can reuse the values below.
export MASTER_KEY="insert the public_key portion of masterkey.json here"
export CONNECTOR_URL=http://localhost:3001/v1
export CLIENT_ID=21jtaatqj8y5t0kctb2ejr6jev5w8
export CLIENT_SECRET=3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw

# Install dependencies and run the sample app.
go get
go run ./app.rb

# In another shell and in the same directory, run grafton.
grafton test --product=bonnets --plan=small --region=aws::us-east-1 \
    --client-id=21jtaatqj8y5t0kctb2ejr6jev5w8 \
    --client-secret=3yTKSiJ6f5V5Bq-kWF0hmdrEUep3m3HKPTcPX7CdBZw \
    --connector-port=3001 \
    --new-plan=large \
    http://localhost:4567

# If everything went well, you'll be greeted with plenty of green check marks!
```
