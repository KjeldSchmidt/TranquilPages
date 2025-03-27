const fs = require('fs');
const path = require('path');

const env = {
  production: false,
  BASE_URL: process.env.BASE_URL,
};

let missing_env_variables = [];
for (const key in env) {
  if (env[key] === undefined) {
    missing_env_variables.push(key);
  }
}

if (missing_env_variables.length) {
  console.error(
      `The following environment variables are not set, but are required for building this app: 
      
      ${missing_env_variables}
      
      Please verify your environment.`
  );
  process.exit(1);
}

const envProd = {
  ...env,
  production: true
};

const environmentsPath = path.join(__dirname, '../src/environments');

fs.writeFileSync(
  path.join(environmentsPath, 'environment.ts'),
  `export const environment = ${JSON.stringify(env, null, 2)};`
);

fs.writeFileSync(
  path.join(environmentsPath, 'environment.prod.ts'),
  `export const environment = ${JSON.stringify(envProd, null, 2)};`
);