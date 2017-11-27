var request = require('request-promise');
const directoryId = '5c37fe2b-b77a-4388-a815-201eb361be16';
const client_id = 'dda51c56-a86e-452f-a91f-cef7b0785526';
const client_secret = 'x5lXipWyEJhXRYV4sCk6Y1+zKzhifa8du7DInNS1lDE=';
const username = 'ublradmin@umbrellarstack2.onmicrosoft.com';
const password = 'AzureStack!';
const resource = 'https://adminmanagement.umbrellarstack2.onmicrosoft.com/c3da154f-8979-424d-a256-3fbd4e9488e2';
const AADURI = "https://login.microsoftonline.com/"+directoryId+"/oauth2/token?api-version=1.0'";
const postHeaders = {"Content-Type": "application/x-www-form-urlencoded"};
process.env['NODE_TLS_REJECT_UNAUTHORIZED'] = 0;

const options = {
    methond : "POST",
    uri: AADURI,
    formData: {
        grant_type: "password",
        scope: "openid",
        resource: resource,
        client_id: client_id,
        client_secret: client_secret,
        username: username,
        password: password
    },
    headers: postHeaders, 
    json: true
}

var accessToken = {};
const suburl = 'https://adminmanagement.asdk2.umbrellar.io/subscriptions/f0dcdd97-386f-4af2-bbbc-26152c89ad07/providers/Microsoft.Subscriptions.Admin/subscriptions?api-version=2015-11-01';

request(options, (err, res, body) => {
        if (err) {
            console.error('Post error: ', err);
        }
        return accessToken = body.access_token;
    }).then (() => {
            request.get(
                suburl, {'auth': {'bearer': accessToken}}, (err, res, body) => {console.log(body)});
    });
            