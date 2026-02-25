# @pittorica/aspect-ratio-react

A component to maintain the aspect ratio of a child element.

## Installation

```bash
npm install @pittorica/aspect-ratio-react
```

You will also need to install the core `pittorica` package which contains the CSS, and its dependencies.

```bash
npm install pittorica @pittorica/box-react
```

## Usage

```jsx
import { AspectRatio } from '@pittorica/aspect-ratio-react';
import { Box } from '@pittorica/box-react';
import 'pittorica/reset';
import 'pittorica/tokens';

function MyImage() {
  return (
    <Box style={{ width: '400px' }}>
      <AspectRatio ratio={16 / 9}>
        <img
          src="https://images.unsplash.com/photo-1769251297393-8178b5988b08?q=80&w=1170&auto=format&fit=crop&ixlib=rb-4.1.0&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D"
          alt="Landscape"
        />
      </AspectRatio>
    </Box>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
