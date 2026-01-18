envmon - Environment file deployment manager

```
Usage:
  envmon <deployment>    Switch to specified deployment configuration
  envmon                 Show current active deployment
  envmon configs         List all detected deployment configurations  
  envmon help            Show this help message
```
## How it works:
  envmon reads the .env file in your current working directory and looks
  for deployment flags in comments. It uncomments lines matching the 
  specified flag while commenting out lines for other flags.

```
Flag Syntax (you can name your flag anything):
  # [staging]            Single deployment flag (active only for staging)
  # [staging/dev]        Multiple flags (active for both staging and dev)
  # [~prod]              Inverse (active for all EXCEPT prod)
  # [~(staging/dev)]     Inverse multiple (active for all EXCEPT staging and dev)


Example (.env file): 
  # [staging] MongoDB config
  #MONGO_URI=staging-uri
  
  # [dev] MongoDB config  
  MONGO_URI=dev-uri

  Running 'envmon staging' will uncomment staging lines and comment dev lines.
```

## Why build a tool?

I built this tool because I often found myself manually commenting/uncommenting multiple sections of my .env during development

This became very time consuming as the env file would be over a thousand lines, and I would often be left doubtful whether I had mixed some of the staging/production keys

Please refer to the `.env.example` file to understand why this tool is needed

### Design Choices

- should be very easy to use, with zero/near to no changes in existing pipeline
- should be very easy to configure, just adding comments before sections in the .env file
- the config itself should be self contained within the env file instead of in a seperate server/database
- while remaining simple in design, allow for writing complex/chained deployment configs

I drew some inspiration from how Vercel and Retool allow for custom deployment environments, seperating env keys across dev, staging and production.

Both are good but are reliant on their web software/infrastructure for key management and only support simple environment seperation.

A good solution should be agnostic of software, infrastructure and frameworks. The developer should have full control over what's active and what's not 