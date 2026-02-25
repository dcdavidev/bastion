# @pittorica/box-react

The `Box` component.

## Installation

```bash
npm install @pittorica/box-react
```

## Usage

```jsx
import { Box } from '@pittorica/box-react';

function MyForm() {
  return (
    <Box as="form">
      <Box as="input" required placeholder="Required field" />
      <Box as="button" type="submit">
        Submit
      </Box>
    </Box>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
