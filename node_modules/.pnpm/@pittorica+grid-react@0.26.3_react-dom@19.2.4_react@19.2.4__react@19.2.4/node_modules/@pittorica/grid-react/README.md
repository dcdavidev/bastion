# @pittorica/grid-react

The `Grid` component.

## Installation

```bash
npm install @pittorica/grid-react
```

## Usage

```jsx
import { Grid, Box } from '@pittorica/grid-react';

function MyDashboard() {
  return (
    <Grid columns={{ initial: '1', md: '2', lg: '4' }} gap="4">
      <Box>Widget 1</Box>
      <Box>Widget 2</Box>
      <Box>Widget 3</Box>
      <Box>Widget 4</Box>
    </Grid>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
