Image Filter API:
----

**Read the below instructions the deploy the application on heroku:**

1. Download and install the gcp config file named *imagefilterapi-config.json* with in the root directory.

2. Setup config in heroku env varibles. Run below commands.

   ```

    heroku config:set GCP_CREDENTIALS="$(< ./imagefilterapi-config.json)" --app image-filtering-for-browser

    heroku config:set GOOGLE_APPLICATION_CREDENTIALS=gcp-config.json --app image-filtering-for-browser

   ```

3. Build the go code

    ```
    GOOS=linux GOARCH=amd64 go build -o ./build/image-filtering-for-browser
    ```

4. Deploy the code on heroku for master branch

5. To check the logs

    ```
    heroku logs --tail --app image-filtering-for-browser
    ```

----

Sample Curl:

----

```

curl -X POST https://image-filtering-for-browser-daf2c3a53445.herokuapp.com/verify_image \
-H "Content-Type: application/json" \
-d '{"image_url": "https://images.unsplash.com/photo-1506748686214-e9df14d4d9d0"}'

```
