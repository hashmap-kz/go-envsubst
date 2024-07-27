# Architecture

- [Project postulated](#project-postulates)
- [Configuration](#configuration)

## Project postulates:

### Simplicity

> **Everywhere: in design, in configuration, in managing**
>
> The reason for that: is you need to be sure what is going on when you start your application.
>
> You have to be sure what this application does, and how to handle all errors as quickly as possible.
>
> It's not a goal to do 1000 things, but here is a main goal: make a few things right.
>
> This project is oriented specifically on Unix based OS, and it will be.

## Configuration

```
# you may use comma or spaces, but not both at the same time

# expand ONLY these list of env-vars and 
# ignore anything else (this is important to not expand nginx conf, for example)
#
export GENVSUBST_ALLOWED='APP_NAME, IMAGE_NAME, IMAGE_TAG'

# do NOT expand these env-vars, never. 
# this setting has the highest priority
#
export GENVSUBST_RESTRICTED='SECRET_X SECRET_Y'

# expand ONLY if the name of env-var starts with any of these prefixes
#
export GENVSUBST_ALLOWED_WITH_PREFIXES='CI_ APP_'

# do NOT expand any of env-var if its name starts with any of these prefixes. 
# this setting is also has the highest priority.
#
export GENVSUBST_RESTRICTED_WITH_PREFIXES='SECRET_ VAULT_'
```

> The combination is as simple as it looks like:
>
> Let's assume, you set all the possible options at the same time.
>
> This pseudocode shows the logic:

```
allowToExpand = (variable not in GENVSUBST_RESTRICTED) 
    && !starts_with(variable, any(GENVSUBST_RESTRICTED_WITH_PREFIXES))
        
if isEmpty(GENVSUBST_ALLOWED) && isEmpty(GENVSUBST_ALLOWED_WITH_PREFIXES) 
{
    if allowToExpand {
        expand(variable)
    }
} else {

    if allowToExpand {
        if (variable in GENVSUBST_ALLOWED) 
            || starts_with(variable, any(GENVSUBST_ALLOWED_WITH_PREFIXES)) {
            expand(variable)
        }
    }    
}
```

### Additional
> You may use internal methods for debug your inputs. The result is look like this:
```
LINE    NAME                             VALUE
6       CI_PROJECT_NAME                  postgres
19      CI_PROJECT_PATH                  cv/system/postgresql
19      CI_COMMIT_REF_NAME               dev
25      CI_PROJECT_NAME                  postgres
38      CI_PROJECT_NAME                  postgres
54      CI_PROJECT_NAME                  postgres
75      CI_PROJECT_ROOT_NAMESPACE        cv
75      CI_PROJECT_NAME                  postgres
85      APP_IMAGE                        postgres:latest
```
> Tokenizer is knows exactly what he should to expand and where.
> 
> Sometimes is suitable to --dry-run and check before injecting the flow in your pipelines.





































