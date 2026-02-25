# @pittorica/skeleton-react

The `Skeleton` component.

## Installation

```bash
npm install @pittorica/skeleton-react
```

## Usage

```jsx
import { Skeleton, Flex, Box } from '@pittorica/skeleton-react';

function MyComponent({ loading }) {
  return (
    <Flex align="center" gap="3">
      <Skeleton loading={loading} variant="circle" width={48} height={48}>
        <img src="avatar.png" alt="User" />
      </Skeleton>
      <Box width={150}>
        <Skeleton loading={loading} variant="text" mb="1" />
        <Skeleton loading={loading} variant="text" width="60%" />
      </Box>
    </Flex>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
