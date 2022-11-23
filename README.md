# GPU node manager

## Google Cloud services

```sh
gcloud services enable compute.googleapis.com iamcredentials.googleapis.com \
    cloudresourcemanager.googleapis.com firestore.googleapis.com
```

### Cloud Firestore

Firestore にアクセスするためのサービス アカウントを作成します。

```sh
export project_id=$( gcloud config get-value project )
gcloud iam service-accounts create firestore-client \
    --display-name "Cloud Firestore Client SA" \
    --description "Service Account for Cloud Firestore client"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:firestore-client@${project_id}.iam.gserviceaccount.com" \
    --role "roles/datastore.user"
gcloud iam service-accounts keys create src/key.json \
    --iam-account "firestore-client@${project_id}.iam.gserviceaccount.com"
```

Firestore にデータベースを作成します。

```sh
gcloud app create --region "asia-northeast1"
gcloud firestore databases create --region "asia-northeast1"
```

### Firebase Authentication

[認証を Firebase で](https://firebase.google.com/docs/auth)行います。

- 認証としてメールアドレス / パスワードを有効化
- プロジェクトの設定から firebase の設定を src/public/js/app/firebase.js に保存します

```sh
firebase.initializeApp({
    apiKey: "your-api-key",
    authDomain: "your-domain.firebaseapp.com",
    projectId: "your-project-id",
    storageBucket: "your-storage-bucket.appspot.com",
    messagingSenderId: "your-message-sender-id",
    appId: "your-app-id"
});
```
