# @pittorica/dialog-react

A window overlaid on either the primary window or another dialog window, rendering the content underneath inert.

## Installation

```bash
npm install @pittorica/dialog-react
```

You will also need to install the core `pittorica` package which contains the CSS, and its dependencies.

```bash
npm install pittorica @pittorica/box-react @pittorica/heading-react @pittorica/button-react
```

## Usage

```jsx
import { useState } from 'react';
import {
  Dialog,
  DialogActions,
  DialogDescription,
  DialogTitle,
} from '@pittorica/dialog-react';
import { Button } from '@pittorica/button-react';
import 'pittorica/reset';
import 'pittorica/tokens';

function MyDialog() {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button variant="elevated" onClick={() => setOpen(true)}>
        Open Dialog
      </Button>

      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>Dialog Title</DialogTitle>
        <DialogDescription>
          This is a description of the dialog content. It provides context for
          the user regarding the action they are about to take.
        </DialogDescription>
        <DialogActions>
          <Button variant="text" color="slate" onClick={() => setOpen(false)}>
            Cancel
          </Button>
          <Button variant="filled" onClick={() => setOpen(false)}>
            Confirm
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
}
```

## API Reference

### Dialog

| Prop                  | Type                             | Description                                      | Default |
| --------------------- | -------------------------------- | ------------------------------------------------ | ------- |
| `open`                | `boolean`                        | If `true`, the dialog is open.                   |         |
| `onClose`             | `() => void`                     | Callback fired when the dialog is closed.        |         |
| `children`            | `React.ReactNode`                | The content of the dialog.                       |         |
| `className`           | `string`                         | Optional class name for the dialog.              |         |
| `appearance`          | `'light' \| 'dark' \| 'inherit'` | The appearance of the dialog.                    |         |
| `closeOnOverlayClick` | `boolean`                        | If `true`, the dialog closes on overlay click.   | `true`  |
| `closeOnEsc`          | `boolean`                        | If `true`, the dialog closes on `Esc` key press. | `true`  |

### DialogTitle

The `DialogTitle` component accepts all props of the `Heading` component.

### DialogDescription

The `DialogDescription` component accepts all props of the `Box` component.

### DialogActions

The `DialogActions` component accepts all props of the `Box` component.
