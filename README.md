# INSTRUCTIONS

## About permission for Firestore db operations from local

1) make sure you already have the service account with the name "App Engine default service account".
2) IAM -> Service Accounts -> for the "App Engine default service account", under Actions, choose "Create Key" -> "JSON" option -> save the json file safely -> save the path to env var GOOGLE_APPLICATION_CREDENTIALS
Note: if you need to create another key, if original key is lost or giving access to another person, the same instructions will work, as multiple keys can be generated and shared.
3) set roles: IAM -> Permissions -> Members -> edit the "App Engine default service account" -> add "Cloud Datastore User" role. or may be "Cloud Datastore Owner" if that is not suffecient -> save
4) run  local.sh and the db operations should work now

## Deploying to cloud

1) appengine:
    "App Engine default service account" should be availible. this (i guess) is created when the appengine component is added to this GCP project.
    otherwise, app deploy to gcp app engine fa ils with below error  
        Updating service [default]...failed.  
        ERROR: (gcloud.app.deploy) Error Response: [13] Error processing user code.  
    the email id of the default service account will be of the below pattern <projectId>@appspot.gserviceaccount.com  
2) "Cloud Build API" api should be enabled for the project. just one time thing.  
3) set the project name correctly in ./scripts/deploy.sh to the `gcloud config set project ` command   
4) run ./scripts/deploy.sh  
