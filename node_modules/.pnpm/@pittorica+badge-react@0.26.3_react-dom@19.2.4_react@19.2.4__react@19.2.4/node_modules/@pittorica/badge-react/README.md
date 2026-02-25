# @pittorica/badge-react

A small numerical value or status descriptor for a UI element.

## Installation

```bash
npm install @pittorica/badge-react
```

You will also need to install the core `pittorica` package which contains the CSS, and its dependencies.

```bash
npm install pittorica @pittorica/box-react @pittorica/text-react
```

## Usage

```jsx
import { Badge } from '@pittorica/badge-react';
import { Box } from '@pittorica/box-react';
import 'pittorica/reset';
import 'pittorica/tokens';

function MyBadge() {
  return (
    <Badge badgeContent={4} color="indigo">
      <Box
        style={{
          width: '40px',
          height: '40px',
          background: 'var(--pittorica-slate-3)',
          borderRadius: '4px',
        }}
      />
    </Badge>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
