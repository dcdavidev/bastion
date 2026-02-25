# @pittorica/avatar-react

An image component with a fallback for representing a user.

## Installation

```bash
npm install @pittorica/avatar-react
```

You will also need to install the core `pittorica` package which contains the CSS, and its dependencies.

```bash
npm install pittorica @pittorica/box-react
```

## Usage

```jsx
import { Avatar } from '@pittorica/avatar-react';
import 'pittorica/reset';
import 'pittorica/tokens';

function MyAvatar() {
  return (
    <Avatar
      src="https://images.unsplash.com/photo-1492633423870-43d1cd2775eb?&w=128&h=128&dpr=2&q=80"
      alt="Colm Tuite"
      fallback="CT"
    />
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
