# CAFE (Configuration Architecture for Flexible Environments)

CAFE is a simple, human-friendly, structured configuration language built with IaaC and configuration files in mind.

The syntax for CAFE is inspired by [Groovy](https://groovy-lang.org/), [NGINX configuration](http://nginx.org/en/docs/beginners_guide.html#conf_structure), and others.

CAFE uses a key-value structure along with hierarchy for better readability.

## Motivations

Why not use something already used, such as YAML, JSON, etc.?

We find that most configuration languages have some downsides that are quite annoying if you’re using it for complex applications. There’s quite a gap between programming languages and markup languages. For example, in a programming language, to declare a variable you can simply go `x = “hello”` (with obvious variations from language to language), but JSON requires all definitions to be inside quotes: `“foo”: “bar”`. Although that is very good for interoperability, it’s annoying if you’re using it to create a configuration file. Another great example is YAML. YAML files can get quite confusing if you need to create a large one, one can easily get lost in the indentation of the file, that is far more annoying that JSON’s quotes.

CAFE attempts to be a bridge between markup / configuration files, and programming languages. It has a syntax similar to what most programming languages look like, but with all the requirements for a markup language. It is made to be easily written and read.

CAFE is build around key-value pairs and a well-defined hierarchy that allows for better readability.

## Syntax

Let’s take the following JSON configuration file for a generic application:

```json
{
  "environment_config": {
    "app": {
      "name": "Sample App",
      "url": "http://10.0.0.120:8080/app/",
    },
    "database": {
      "name": "mysql database",
      "host": "10.0.0.120",
      "port": 3128,
      "username": "root",
      "password": "toor",
    },
    "rest_api": "http://10.0.0.120:8080/v2/api/"
  }
}
```

The CAFE equivalent of this configuration is the following:

```
environment_config {
    app {
        name = "Sample App"
        url = "http://10.0.0.120:8080/app/"
    }
    database {
      name = "mysql database"
      host = "10.0.0.120"
      port = 3128
      username = "root"
      password = "toor"
    }
    rest_api = "http://10.0.0.120:8080/v2/api/"
}
```

We can convert other languages too, for example, this Kubernetes deployment file (YAML):

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
```

Has the following CAFE equivalent:

```
globalAppName = "nginx"
apiVersion = "apps/v1"
kind = "Deployment"
metadata {
    name = "${globalAppName}-deployment"
    labels {
        app = globalAppName
    }
}
spec {
    replicas = 3
    selector {
        matchLabels {
            app = globalAppName
        }
    }
    template {
        metadata {
            labels {
                app = globalAppName
            }
        }
        spec {
            containers {
                nginx {
                    image = "nginx:1.14.2"
                    ports {
                        containerPort = 80
                    }
                }
            }
        }
    }
}
```

As you can see, CAFE can be used in many ways, with better readability and usage.

CAFE also supports expressions:

```
// String interpolations
text = "world"
message = "Hello, ${text}!"

// Arithmetic operations
v1 = 1
v2 = 1
sum = v1 + v2

// Comparisons
check_sum = v1 == 1

// Functions
up = upper(message)
```

For more information, check the [syntax spec document](SPEC.md)

## References

- [Why JSON isn’t a Good Configuration Language](https://www.lucidchart.com/techblog/2018/07/16/why-json-isnt-a-good-configuration-language/)
- [Don’t Use JSON as a Configuration File Format. (Unless Absolutely You Have To…)](https://revelry.co/insights/development/json-configuration-file-format/)
- [Why YAML is used for configuration when it's so bad and what can you do about it?](https://kula.blog/posts/yaml/)
- [The state of config file formats: XML vs. YAML vs. JSON vs. HCL](https://octopus.com/blog/state-of-config-file-formats)