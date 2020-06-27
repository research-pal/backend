




About permission for Firestore db operations from local:
1) make sure you already have the service account with the name "App Engine default service account".
2) IAM -> Service Accounts -> for the "App Engine default service account", under Actions, choose "Create Key" -> "JSON" option -> save the json file safely -> save the path to env var GOOGLE_APPLICATION_CREDENTIALS
3) set roles: IAM -> Permissions -> Members -> edit the "App Engine default service account" -> add "Cloud Datastore User" role. or may be "Cloud Datastore Owner" if that is not suffecient -> save
4) run  local.sh and the db operations should work now

