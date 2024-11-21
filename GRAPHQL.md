```gql
query Device($device_name: String!) {
  InfraDevice(name__value: $device_name) {
    edges {
      node {
        id
        name {
          value
        }
        role {
          value
        }
        platform {
          node {
            id
          }
        }
        primary_address {
          node {
            id
          }
        }
        status {
          id
        }
        topology {
          node {
            id
          }
        }
        device_type {
          node {
            id
          }
        }
        asn {
          node {
            asn {
              id
            }
          }
        }
      }
    }
  }
}
```