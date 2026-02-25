# @pittorica/toast-react

The `Toast` component.

## Installation

```bash
npm install @pittorica/toast-react
```

## Usage

```jsx
import { toast, ToastProvider } from '@pittorica/toast-react';

function MyApp() {
  return (
    <>
      <ToastProvider />
      <button
        onClick={() =>
          toast({
            title: 'Success!',
            description: 'Action completed.',
            color: 'teal',
          })
        }
      >
        Show Toast
      </button>
    </>
  );
}
```

## License

This project is licensed under the MIT License.

**Copyright (c) 2025 Davide Di Criscito**

For the full details, see the [LICENSE](LICENSE) file.
