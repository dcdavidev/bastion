# @pittorica/inset-react

The `Inset` component.

## Installation

```bash
npm install @pittorica/inset-react
```

## Usage

```jsx
import { Inset, Box, Text } from '@pittorica/react';

function MyCard() {
  return (
    <Box p="4" style={{ border: '1px solid var(--pittorica-slate-4)' }}>
      <Inset side="top">
        <img src="header.jpg" alt="Header" />
      </Inset>
      <Text mt="3">Card Content</Text>
    </Box>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
