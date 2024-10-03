pipeline {
    agent any
    stages {
        stage('Make HTTP Request') {
            steps {
                script {
                    // HTTP request using all supported options
                    def response = httpRequest(
                        url: 'https://jsonplaceholder.typicode.com/posts',  // URL to send the request to
                        httpMode: 'POST',                                   // HTTP method (GET, POST, PUT, DELETE, etc.)
                        authentication: 'httprequest',                      // Jenkins credentials ID for authentication
                        contentType: 'APPLICATION_JSON',                    // Content-Type of the request
                        acceptType: 'APPLICATION_JSON',                     // Accept-Type of the response
                        customHeaders: [[name: 'Authorization', value: 'Bearer <token>'], [name: 'X-Custom-Header', value: 'example-header']],  // Custom headers
                        requestBody: '{"title": "foo", "body": "bar", "userId": 1}',  // Body for POST requests
                        validResponseCodes: '200:299',                      // Range of valid response codes
                        validResponseContent: '"id":',                      // Expected content in the response body
                        outputFile: 'response.json',                        // Save the response body to a file
                        timeout: 60,                                        // Timeout in seconds
                        ignoreSslErrors: true,                              // Ignore SSL certificate errors
                        consoleLogResponseBody: true,                       // Log the response body in the console
                        quiet: false,                                       // Log detailed information
                        wrapAsMultipart: false                              // Send the request without wrapping as multipart
                    )

                    // Access the response details
                    echo "Response Code: ${response.status}"
                    echo "Response Content: ${response.content}"
                    echo "Response Headers: ${response.headers}"
                }
            }
        }
        stage('Read Response') {
            steps {
                // Read and print the content of the saved response file
                script {
                    def responseContent = readFile 'response.json'
                    echo "Downloaded Response: ${responseContent}"
                }
            }
        }
    }
}