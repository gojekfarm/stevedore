# **Plan**

###### **With Specific Files:**

`stevedore apply -f example.yaml -o example_override_example.yaml -e dev_environment.yaml`

###### **With directories:**
`stevedore apply -f services/ -o overrides/ -e envs/`

###### **What it does:**

Stevedore shows the environment information to the user and waits for confirmation from the user to apply the changes shown by the [plan](plan.md) command. 

Upon confirmation stevedore applies the operation. The confirmation step can be skipped by providing the `--yes` flag

**Flags accepted:**

|  Flags 	| Description  	| Notes  	|
|---	|---	|---	|
| -f  	|  manifest file or directory with manifest files 	| Atleast one should be given |
| -o  	| override file or directory with override files  	| Optional |
| -e  	| env file or directory with env files  	| Optional |
| -a  	| output of artifacts path 	| Optional |
| -r  	| helm repo name to which the charts will be pushed/pulled  	| Optional. Defaults to `chartmuseum` |
| -t  	| helm timeout  	| Optional. Defaults to `600 seconds` |
| --yes  	| skips the confirmation steps and applies the operation  	| Optional |
| --recursive  	| respects `depends_on` in manifest files and recursively renders/plans/applies for those manifests too  	| Optional|
| --config  	|  config file that stevedore looks at for environment information 	| Optional. Defaults to `/usr/local/etc/stevedore-config.yaml` |
| --log-level  	|  sets the logger level 	| Optional. Defaults to `error` |


_**Atleast one manifest file should be given**_. Overrides/Envs file(s) are optional. If provided they will be used.

Stevedore goes through all the manifests, collates all matching overrides and envs for the appropriate manifests.
It shows the various kubernetes entities created/modified/deleted. 

A very useful detail to be noted is the `depends_on` portion. Using this, one manifest can depend on another manifest.
If the recursive flag is provided, now stevedore will recursively also do a render/plan/apply for the manifests dependencies mentioned in this part of the manifest. 

This command ideally should precede the [apply](apply.md)  command. 

Todo:

Add gif