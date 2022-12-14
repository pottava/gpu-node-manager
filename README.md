# GPU node manager

## リポジトリのクローン

```sh
git clone https://github.com/pottava/gpu-node-manager.git
cd gpu-node-manager/
```

## Google Cloud services

```sh
gcloud services enable compute.googleapis.com iamcredentials.googleapis.com \
    cloudresourcemanager.googleapis.com firestore.googleapis.com \
    cloudbuild.googleapis.com appengine.googleapis.com run.googleapis.com \
    artifactregistry.googleapis.com containerscanning.googleapis.com \
    notebooks.googleapis.com aiplatform.googleapis.com \
    secretmanager.googleapis.com
```

### Cloud IAM

アプリケーションのためのサービス アカウントを作成します。

```sh
export project_id=$( gcloud config get-value project )
gcloud iam service-accounts create app-client \
    --display-name "SA for the app" \
    --description "Service Account for the GPU application"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:app-client@${project_id}.iam.gserviceaccount.com" \
    --role "roles/notebooks.admin"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:app-client@${project_id}.iam.gserviceaccount.com" \
    --role "roles/datastore.user"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:app-client@${project_id}.iam.gserviceaccount.com" \
    --role "roles/resourcemanager.projectIamAdmin"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:app-client@${project_id}.iam.gserviceaccount.com" \
    --role "roles/storage.admin"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:app-client@${project_id}.iam.gserviceaccount.com" \
    --role "roles/firebaseauth.admin"
```

### Cloud Firestore

Firestore にデータベースを作成します。

```sh
gcloud app create --region "asia-northeast1"
gcloud firestore databases create --region "asia-northeast1"
```

### Firebase Authentication

[認証を Firebase で](https://firebase.google.com/docs/auth)行います。

https://console.firebase.google.com/

1. "Authentication" の認証、もしくは [Identity Platform: ID プロバイダ](https://console.cloud.google.com/customer-identity/providers) からメールアドレス / パスワードを有効化
2. "Authentication" の Users、もしくは [Identity Platform: ユーザー](https://console.cloud.google.com/customer-identity/users) からユーザーを登録
3. "プロジェクトの設定" で確認できる内容を src/public/js/app/firebase.js に保存

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

### Cloud Run

公開 URL を取得するため、サンプルアプリケーションをデプロイしておきます。

```sh
gcloud run deploy dev --image gcr.io/cloudrun/hello --region "asia-northeast1" \
    --platform "managed" --cpu 1.0 --memory 128Mi --max-instances 2 \
    --allow-unauthenticated
gcloud run services add-iam-policy-binding dev --region "asia-northeast1" \
    --member "allUsers" --role "roles/run.invoker"
gcloud run deploy prod --image gcr.io/cloudrun/hello --region "asia-northeast1" \
    --platform "managed" --cpu 1.0 --memory 128Mi --max-instances 10 \
    --allow-unauthenticated
gcloud run services add-iam-policy-binding prod --region "asia-northeast1" \
    --member "allUsers" --role "roles/run.invoker"
```

## ローカル開発

```sh
go install github.com/revel/cmd/revel@latest
export GOOGLE_CLOUD_PROJECT=$( gcloud config get-value project )
gcloud auth application-default login
revel run -a src
```

## CI / CD パイプライン

以下を設定後、git push により開発環境の Cloud Run が更新されることを確認してください。

### GitHub Actions

CI ツールに渡すサービスアカウントを作ります。

```sh
export project_id=$( gcloud config get-value project )
gcloud iam service-accounts create cicd-service \
    --display-name "SA for CI/CD" \
    --description "Service Account for CI/CD pipelines"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:cicd-service@${project_id}.iam.gserviceaccount.com" \
    --role "roles/run.admin"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:cicd-service@${project_id}.iam.gserviceaccount.com" \
    --role "roles/storage.admin"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:cicd-service@${project_id}.iam.gserviceaccount.com" \
    --role "roles/artifactregistry.writer"
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:cicd-service@${project_id}.iam.gserviceaccount.com" \
    --role "roles/cloudbuild.builds.editor"
gcloud iam service-accounts add-iam-policy-binding \
    app-client@${project_id}.iam.gserviceaccount.com \
    --member "serviceAccount:cicd-service@${project_id}.iam.gserviceaccount.com" \
    --role "roles/iam.serviceAccountUser"
gcloud iam service-accounts keys create key.json \
    --iam-account "cicd-service@${project_id}.iam.gserviceaccount.com"
cat key.json && rm -f key.json
```

GitHub プロジェクトの Secret に以下の値を設定します。

- GOOGLECLOUD_PROJECT: プロジェクト ID
- GOOGLECLOUD_SA_KEY: デプロイするためのサービス アカウント
- GOOGLECLOUD_FIREBASE: Firebase の設定 JSON（ダブル クオーテーションにはエスケープが必要）

### Cloud Build

Artifact Registry にリポジトリを作成します。

```sh
gcloud artifacts repositories create gpu-node-manager \
    --repository-format docker --location asia-northeast1 \
    --description="Docker repository for GPU node manager"
```

Secret Manager に Firebase の設定を保存します。

```sh
gcloud secrets create firebase-configs --replication-policy "automatic" \
    --data-file src/public/js/app/firebase.js
```

Cloud Build に必要となる権限を付与します。

```sh
export project_id=$( gcloud config get-value project )
export project_number=$(gcloud projects describe ${project_id} \
    --format="value(projectNumber)")
gcloud projects add-iam-policy-binding "${project_id}" \
    --member "serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" \
    --role "roles/run.developer"
gcloud secrets add-iam-policy-binding firebase-configs \
    --member "serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" \
    --role roles/secretmanager.secretAccessor
gcloud iam service-accounts add-iam-policy-binding \
    app-client@${project_id}.iam.gserviceaccount.com \
    --member "serviceAccount:${project_number}@cloudbuild.gserviceaccount.com" \
    --role "roles/iam.serviceAccountUser"
```
