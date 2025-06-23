# Kubernetes Manifests for Hackz Giganoto

このディレクトリには、Hackz GiganoToアプリケーションをKubernetesにデプロイするためのマニフェストファイルが含まれています。

## 構成

```
k8s-manifests/
├── base/                    # ベースマニフェスト
│   ├── kustomization.yaml
│   ├── namespace.yaml
│   ├── configmaps/         # 設定ファイル
│   └── services/           # 各サービスのマニフェスト
├── secrets/                # 秘密情報（GitHub OAuth、JWT等）
├── overlays/               # 環境別設定
│   ├── development/        # 開発環境用
│   └── production/         # 本番環境用
└── README.md
```

## デプロイ方法

### 前提条件

1. Kubernetes クラスターが利用可能
2. `kubectl` と `kustomize` がインストール済み
3. 必要な秘密情報が `secrets/` ディレクトリに設定済み

### 開発環境へのデプロイ

```bash
# 開発環境用名前空間の作成
kubectl create namespace hackz-giganoto-dev

# 開発環境へのデプロイ
kubectl apply -k k8s-manifests/overlays/development
```

### 本番環境へのデプロイ

```bash
# 本番環境用名前空間の作成
kubectl create namespace hackz-giganoto

# 本番環境へのデプロイ
kubectl apply -k k8s-manifests/overlays/production
```

### ArgoCD を使用したCD

ArgoCD を使用してこのマニフェストをデプロイする場合：

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: hackz-giganoto-dev
  namespace: argocd
spec:
  project: default
  source:
    repoURL: https://github.com/naoido/hackz-giganoto
    targetRevision: main
    path: k8s-manifests/overlays/development
  destination:
    server: https://kubernetes.default.svc
    namespace: hackz-giganoto-dev
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

## 秘密情報の設定

デプロイ前に以下の秘密情報を適切に設定してください：

### Auth Secrets (`secrets/auth-secrets.yaml`)

```bash
# GitHub OAuth設定
kubectl create secret generic auth-secrets \
  --from-literal=github-client-id="your-client-id" \
  --from-literal=github-client-secret="your-client-secret" \
  --from-literal=github-redirect-url="https://your-domain/auth/github/callback" \
  -n hackz-giganoto
```

### JWT Secret (`secrets/jwt-secret.yaml`)

```bash
# JWT署名用秘密鍵
kubectl create secret generic jwt-secret \
  --from-literal=jwt-secret-key="your-jwt-secret-key-32-chars-min" \
  -n hackz-giganoto
```

### Grafana Secrets (`secrets/grafana-secrets.yaml`)

```bash
# Grafana管理者認証
kubectl create secret generic grafana-secrets \
  --from-literal=admin-user="admin" \
  --from-literal=admin-password="your-secure-password" \
  -n hackz-giganoto
```

## サービス構成

### マイクロサービス

- **auth**: 認証サービス (HTTP:8000)
- **chat**: チャットサービス (gRPC:50053)
- **bff**: Backend for Frontend (gRPC:50054)

### インフラサービス

- **kong**: APIゲートウェイ (HTTP:8000, Admin:8001)
- **redis**: 各マイクロサービス専用
- **otel-collector**: OpenTelemetryコレクター
- **jaeger**: 分散トレーシング (UI:16686)
- **loki**: ログ集約 (API:3100)
- **grafana**: 監視UI (UI:3000)

## CI/CD パイプライン

GitHub Actions を使用したCI/CDパイプラインが設定されています：

### ビルド・プッシュ (`.github/workflows/build-and-push.yml`)

- リリース作成時に自動実行
- 各マイクロサービスのDockerイメージをビルド
- GitHub Container Registry (ghcr.io) にプッシュ
- Kustomizationのイメージタグを自動更新

### テスト (`.github/workflows/test.yml`)

- PR・プッシュ時に実行
- Go テスト、リント、セキュリティスキャン
- Kubernetesマニフェストの検証

## トラブルシューティング

### ポッド起動の確認

```bash
# 全ポッドの状態確認
kubectl get pods -n hackz-giganoto

# ログ確認
kubectl logs -f deployment/auth-deployment -n hackz-giganoto
```

### サービス接続の確認

```bash
# サービス一覧
kubectl get svc -n hackz-giganoto

# Kong Gateway経由でのヘルスチェック
kubectl port-forward svc/kong-service 8000:8000 -n hackz-giganoto
curl http://localhost:8000/auth/health
```

### 設定の確認

```bash
# ConfigMap確認
kubectl get configmap -n hackz-giganoto
kubectl describe configmap kong-config -n hackz-giganoto

# Secret確認
kubectl get secrets -n hackz-giganoto
```

## 監視・ログ

- **Grafana**: `http://grafana-service:3000` (admin/admin)
- **Jaeger**: `http://jaeger-service:16686`
- **Kong Admin**: `http://kong-service:8001`

## セキュリティ注意事項

1. 本番環境では適切な秘密情報を設定してください
2. LoadBalancer のソース IP 制限を適切に設定してください
3. RBAC を適切に設定してください
4. Network Policy を使用してポッド間通信を制限することを検討してください