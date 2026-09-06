# Working with SparkConnect

A `SparkConnect` is a `v1alpha1` custom resource that represents a long-running [Spark Connect](https://spark.apache.org/docs/latest/spark-connect-overview.html) server (the driver) and its dynamically-managed executor pods.

## Creating a SparkConnect

A `SparkConnect` can be created from a YAML file using `kubectl apply -f <file>`. The operator creates the server pod, the supporting `Service`, and a `ConfigMap` containing the executor pod template. The executor pods are created by Spark based on the configuration in the server pod's `spark-submit` arguments.

```yaml
apiVersion: sparkoperator.k8s.io/v1alpha1
kind: SparkConnect
metadata:
  name: spark-connect
  namespace: default
spec:
  sparkVersion: 4.0.4
  image: "docker.io/apache/spark:4.0.4"
  server:
    cores: 1
    coreRequest: "500m"
    coreLimit: "1"
    memory: 1g
  executor:
    instances: 2
    cores: 1
    coreRequest: "500m"
    coreLimit: "1500m"
    memory: 512m
```

## Deleting a SparkConnect

A `SparkConnect` can be deleted with `kubectl delete sparkconnect <name>`. Deleting the `SparkConnect` causes the operator to garbage-collect the server pod, the `Service`, and the `ConfigMap`. Executor pods are owned by Spark and are removed when the server shuts down.

## Updating a SparkConnect

A `SparkConnect` can be updated using `kubectl apply -f <updated file>`. The admission webhook re-validates the spec on every update.

The operator builds the server pod only on first creation, so changes to `spec.server.coreRequest`, `spec.server.coreLimit`, `spec.server.cores`, `spec.server.memory`, `spec.server.template`, or `spec.image` only take effect after the server pod is recreated (delete the existing pod, or delete and re-apply the `SparkConnect`). Other fields such as `spec.executor.*`, `spec.sparkConf`, `spec.dynamicAllocation`, and `spec.server.service` are applied to subsequent pod creations and Spark submissions.

## Specifying CPU Resources

`SparkPodSpec` exposes `cores` (Spark task-slot count, mapped to `spark.driver.cores` or `spark.executor.cores`) and `coreRequest` / `coreLimit` (physical Kubernetes CPU request/limit, mapped to the container's `resources.{requests,limits}.cpu`). `cores` and `coreRequest` / `coreLimit` are independent.

The server pod is created directly by the operator, so `spec.server.coreRequest` / `coreLimit` are applied to the operator-created server pod's container resources. The executor pods are created by Spark, so `spec.executor.coreRequest` / `coreLimit` are passed via `spark.kubernetes.executor.{request,limit}.cores`. The admission webhook enforces a positive value for both fields and `coreRequest <= coreLimit` when both are set.

## Connecting a Spark Client

Once the server pod is ready, the operator exposes it via a `<name>-server` `Service` on port `15002`:

```bash
${SPARK_HOME}/bin/spark-shell --remote sc://spark-connect-server.default.svc:15002
```

## Dynamic Allocation

Dynamic allocation is configured via `spec.dynamicAllocation` and requires Spark 3.0+. It lets Spark scale the executor pool based on workload. See the [Apache Spark documentation](https://spark.apache.org/docs/latest/configuration.html#dynamic-allocation) for semantics.

## Field Reference

For the full API definition, see the [`SparkConnect` API reference](../reference/api-docs.md).
