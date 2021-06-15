# Push your /wiki folder to Github Wiki

I wrote an overly complicated Go Action to update you github wiki with the content in the /wiki folder of your repo.

You can use a sample action file like this:

```yaml
name: Documentation

on:
  push:
    paths:
      # Trigger only when wiki directory changes
      - "wiki/**"
    branches:
      # And only on primary branch
      - master

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v1
      # Additional steps to generate documentation in "Documentation" directory
      - name: Upload Documentation to Wiki
        uses: eohlde/push_to_wiki@v0
        env:
          GH_PERSONAL_ACCESS_TOKEN: ${{ secrets.GH_PERSONAL_ACCESS_TOKEN }}
```

## Personal Access Token

You will have to create your own personal access token and create a registry secret.
