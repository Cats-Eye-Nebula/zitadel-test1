[![Netlify Status](https://api.netlify.com/api/v1/badges/b82a23f7-d8c7-4025-af18-a46586e89ed0/deploy-status)](https://app.netlify.com/sites/zitadel-docs/deploys)

# Website

This website is built using [Docusaurus 2](https://v2.docusaurus.io/), a modern static website generator.

## Add new Sites to existing Topics

To add a new site to the already existing structure simply save the `md` file into the corresponding folder and append the sites id int the file `sidebars.js`.

## Installation

```console
yarn install
```

## Local Development

```console
yarn start
```

As an alternative you can use this docker compose command:

```console
docker compose up
```

These commands start a local development server.
Most changes are reflected live without having to restart the server.

## Build

```console
yarn build
```

This command generates static content into the `build` directory and can be served using any static contents hosting service.

## Deployment

Each PR will be automatically built with a preview link from cloudflare pages.
Visit the checks / comments on the PR for the link.
