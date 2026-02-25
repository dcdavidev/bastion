# @pittorica/progress-react

The `Progress` component.

## Installation

```bash
npm install @pittorica/progress-react
```

## Usage

```jsx
import { Progress } from '@pittorica/progress-react';

function MyComponent() {
  return (
    <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
      <Progress value={45} color="indigo" />
      <Progress value={100} color="success" />
      <Progress value={15} color="danger" />
    </div>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
