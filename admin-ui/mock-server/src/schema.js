const { gql } = require('apollo-server');

const typeDefs = gql`
  type Query {
    me: User
    runtimes: [Runtime]!
    versions(runtimeId: ID!): [Version]!
    domains: [Domain]!
    usersActivity: [UserActivity]!
    runtime(id: ID!): Runtime
    settings: Settings!
    dashboard: Dashboard!
  }

  type Mutation {
    createRuntime(name: String!): RuntimeUpdateResponse!
    removeRuntime(id: ID!): RuntimeUpdateResponse!
    setSettings(input: SettingsInput): Settings
    addAllowedDomain(domainName: String!): Settings
    removeAllowedDomain(domainName: String!): Settings
    uploadVersion(name: String, type: String!, file: Upload!): KrtFile!
  }

  type User {
    id: ID!
    email: String!
    disabled: Boolean
  }

  type Settings {
    authAllowedDomains: [String]!
    cookieExpirationTime: Int!
  }
  input SettingsInput {
    authAllowedDomains: [String]
    cookieExpirationTime: Int
  }

  type Dashboard {
    runtimes: [Runtime]
    alerts: [Alert]
  }

  type Runtime {
    id: ID!
    name: String!
    status: String!
    creationDate: String!
    versions(status: String!): [Version]
  }

  type Version {
    id: ID!
    versionNumber: String!
    description: String
    status: String!
    creationDate: String!
    creatorName: String!
    activationDate: String!
    activatorName: String!
  }

  type KrtFile {
    filename: String!
    mimetype: String!
    encoding: String!
  }

  type Alert {
    id: ID!
    type: String!
    message: String!
    runtime: Runtime!
  }

  type Domain {
    id: ID!
    name: String!
  }

  type UserActivity {
    id: ID!
    user: String!
    message: String!
    date: String!
  }

  type RuntimeUpdateResponse {
    success: Boolean!
    message: String
    runtime: Runtime!
  }

  enum RuntimeStatus {
    CREATED
    ACTIVE
    WARNING
    ERROR
  }
  enum AlertLevel {
    ERROR
    WARNING
  }
  enum VersionType {
    FIX
    MINOR_UPDATE
    MAJOR_UPDATE
  }
`;

module.exports = typeDefs;
