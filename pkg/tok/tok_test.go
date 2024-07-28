package tok

import (
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

func mustTokenizeString(b string) *Tokenlist {
	tl, _ := Tokenize(b)
	return tl
}

func TestUnexpanded(t *testing.T) {
	input := `
# TODO: MASTER
# ---
# apiVersion: external-secrets.io/v1beta1
# kind: ExternalSecret
# metadata:
#   name: &app ${CI_PROJECT_NAME}
# spec:
#   refreshInterval: "15s"
#   secretStoreRef:
#     name: cluster-secret-store
#     kind: ClusterSecretStore
#   target:
#     name: *app
#     template:
#       type: Opaque
#       engineVersion: v2
#   dataFrom:
#     - extract:
#         key: secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME}

---
apiVersion: v1
kind: Secret
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
stringData:
  POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=ru_RU.UTF-8 --lc-ctype=ru_RU.UTF-8"
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: postgres
  POSTGRES_DB: postgres

---
apiVersion: v1
kind: Service
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
spec:
  type: NodePort
  selector:
    app: *app
  ports:
    - port: 5432
      # dev/tst/preprod/master
      nodePort: 31200

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
spec:
  replicas: 1
  serviceName: *app
  selector:
    matchLabels:
      app: *app
  template:
    metadata:
      labels:
        app: *app
    spec:
      # dev/tst/preprod/master
      nodeSelector:
        worker5: worker5

      volumes:
        - name: storage
          persistentVolumeClaim:
            claimName: ${CI_PROJECT_ROOT_NAMESPACE}-${CI_PROJECT_NAME}-data
        - name: dshm
          emptyDir:
            medium: Memory

      imagePullSecrets:
        - name: myprivateregistry

      containers:
        - name: *app
          image: ${APP_IMAGE}
          imagePullPolicy: Always
          terminationMessagePolicy: FallbackToLogsOnError
          env:
            - name: TZ
              value: Asia/Aqtau
          envFrom:
            - secretRef:
                name: *app
          # dev/tst/preprod/master
          resources:
            limits:
              cpu: "2"
              memory: 2Gi
            requests:
              cpu: 128m
              memory: 256Mi
          volumeMounts:
            - mountPath: /var/lib/postgresql/data
              name: storage
            - mountPath: /dev/shm
              name: dshm
`
	input = strings.TrimSpace(input)

	tokenlist := mustTokenizeString(input)
	tokenized := tokenlist.DumpRawUnexpanded()

	assert.Equal(t, input, tokenized)
}

func TestExpanded(t *testing.T) {
	input := `
# TODO: MASTER
# ---
# apiVersion: external-secrets.io/v1beta1
# kind: ExternalSecret
# metadata:
#   name: &app ${CI_PROJECT_NAME}
# spec:
#   refreshInterval: "15s"
#   secretStoreRef:
#     name: cluster-secret-store
#     kind: ClusterSecretStore
#   target:
#     name: *app
#     template:
#       type: Opaque
#       engineVersion: v2
#   dataFrom:
#     - extract:
#         key: secret/${CI_PROJECT_PATH}/${CI_COMMIT_REF_NAME}

---
apiVersion: v1
kind: Secret
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
stringData:
  POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=ru_RU.UTF-8 --lc-ctype=ru_RU.UTF-8"
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: postgres
  POSTGRES_DB: postgres

---
apiVersion: v1
kind: Service
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
spec:
  type: NodePort
  selector:
    app: *app
  ports:
    - port: 5432
      # dev/tst/preprod/master
      nodePort: 31200

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: &app ${CI_PROJECT_NAME}
  labels:
    app: *app
spec:
  replicas: 1
  serviceName: *app
  selector:
    matchLabels:
      app: *app
  template:
    metadata:
      labels:
        app: *app
    spec:
      # dev/tst/preprod/master
      nodeSelector:
        worker5: worker5

      volumes:
        - name: storage
          persistentVolumeClaim:
            claimName: ${CI_PROJECT_ROOT_NAMESPACE}-${CI_PROJECT_NAME}-data
        - name: dshm
          emptyDir:
            medium: Memory

      imagePullSecrets:
        - name: myprivateregistry

      containers:
        - name: *app
          image: ${APP_IMAGE}
          imagePullPolicy: Always
          terminationMessagePolicy: FallbackToLogsOnError
          env:
            - name: TZ
              value: Asia/Aqtau
          envFrom:
            - secretRef:
                name: *app
          # dev/tst/preprod/master
          resources:
            limits:
              cpu: "2"
              memory: 2Gi
            requests:
              cpu: 128m
              memory: 256Mi
          volumeMounts:
            - mountPath: /var/lib/postgresql/data
              name: storage
            - mountPath: /dev/shm
              name: dshm
`
	expected := `
# TODO: MASTER
# ---
# apiVersion: external-secrets.io/v1beta1
# kind: ExternalSecret
# metadata:
#   name: &app postgres
# spec:
#   refreshInterval: "15s"
#   secretStoreRef:
#     name: cluster-secret-store
#     kind: ClusterSecretStore
#   target:
#     name: *app
#     template:
#       type: Opaque
#       engineVersion: v2
#   dataFrom:
#     - extract:
#         key: secret/cv/system/postgresql/dev

---
apiVersion: v1
kind: Secret
metadata:
  name: &app postgres
  labels:
    app: *app
stringData:
  POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --lc-collate=ru_RU.UTF-8 --lc-ctype=ru_RU.UTF-8"
  POSTGRES_USER: postgres
  POSTGRES_PASSWORD: postgres
  POSTGRES_DB: postgres

---
apiVersion: v1
kind: Service
metadata:
  name: &app postgres
  labels:
    app: *app
spec:
  type: NodePort
  selector:
    app: *app
  ports:
    - port: 5432
      # dev/tst/preprod/master
      nodePort: 31200

---
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: &app postgres
  labels:
    app: *app
spec:
  replicas: 1
  serviceName: *app
  selector:
    matchLabels:
      app: *app
  template:
    metadata:
      labels:
        app: *app
    spec:
      # dev/tst/preprod/master
      nodeSelector:
        worker5: worker5

      volumes:
        - name: storage
          persistentVolumeClaim:
            claimName: cv-postgres-data
        - name: dshm
          emptyDir:
            medium: Memory

      imagePullSecrets:
        - name: myprivateregistry

      containers:
        - name: *app
          image: postgres:latest
          imagePullPolicy: Always
          terminationMessagePolicy: FallbackToLogsOnError
          env:
            - name: TZ
              value: Asia/Aqtau
          envFrom:
            - secretRef:
                name: *app
          # dev/tst/preprod/master
          resources:
            limits:
              cpu: "2"
              memory: 2Gi
            requests:
              cpu: 128m
              memory: 256Mi
          volumeMounts:
            - mountPath: /var/lib/postgresql/data
              name: storage
            - mountPath: /dev/shm
              name: dshm`

	os.Setenv("CI_PROJECT_ROOT_NAMESPACE", "cv")
	os.Setenv("CI_PROJECT_PATH", "cv/system/postgresql")
	os.Setenv("CI_PROJECT_NAME", "postgres")
	os.Setenv("CI_COMMIT_REF_NAME", "dev")
	os.Setenv("APP_IMAGE", "postgres:latest")

	input = strings.TrimSpace(input)
	expected = strings.TrimSpace(expected)

	tokenlist := mustTokenizeString(input)
	tokenized := tokenlist.DumpExpanded()

	assert.Equal(t, expected, tokenized)
}

func TestRestricted(t *testing.T) {
	os.Setenv("SECRET", "1")
	os.Setenv("XXX", "2")
	os.Setenv("GENVSUBST_RESTRICTED", "SECRET")

	input := "${SECRET}-${XXX}-${SECRET}"
	tokenlist := mustTokenizeString(input)
	tokenized := tokenlist.DumpExpanded()

	assert.Equal(t, "${SECRET}-2-${SECRET}", tokenized)
}

func TestRestrictedWithPrefixes(t *testing.T) {
	os.Setenv("SECRET", "1")
	os.Setenv("XXX", "2")
	os.Setenv("GENVSUBST_RESTRICTED_WITH_PREFIXES", "SECR")

	input := "${SECRET}-${XXX}-${SECRET}"
	tokenlist := mustTokenizeString(input)
	tokenized := tokenlist.DumpExpanded()

	assert.Equal(t, "${SECRET}-2-${SECRET}", tokenized)
}

func TestAllowed(t *testing.T) {
	os.Setenv("var1", "1")
	os.Setenv("var2", "2")
	os.Setenv("var3", "3")
	os.Setenv("GENVSUBST_ALLOWED", "var2")

	input := "${var1}-${var2}-${var3}"
	tokenlist := mustTokenizeString(input)
	tokenized := tokenlist.DumpExpanded()

	assert.Equal(t, "${var1}-2-${var3}", tokenized)
}

func TestAllowedWithPrefixes(t *testing.T) {
	os.Setenv("var1", "1")
	os.Setenv("var2", "2")
	os.Setenv("var3", "3")
	os.Setenv("ANOTHER", "@@@")
	os.Setenv("GENVSUBST_ALLOWED_WITH_PREFIXES", "var")

	input := "${var1}-${var2}-${var3}-${ANOTHER}"
	tokenlist := mustTokenizeString(input)
	tokenized := tokenlist.DumpExpanded()

	assert.Equal(t, "1-2-3-${ANOTHER}", tokenized)
}
