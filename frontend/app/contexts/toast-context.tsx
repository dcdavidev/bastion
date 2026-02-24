import { createContext, useContext, useState, ReactNode } from "react";
import { Box, Text, Stack } from "@pittorica/react";
import { X, CheckCircle, AlertCircle } from "lucide-react";

type ToastType = "success" | "error" | "info";

interface Toast {
  id: string;
  message: string;
  type: ToastType;
}

interface ToastContextType {
  toast: (message: string, type?: ToastType) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const toast = (message: string, type: ToastType = "info") => {
    const id = Math.random().toString(36).substring(2, 9);
    setToasts((prev) => [...prev, { id, message, type }]);
    setTimeout(() => {
      setToasts((prev) => prev.filter((t) => t.id !== id));
    }, 5000);
  };

  return (
    <ToastContext.Provider value={{ toast }}>
      {children}
      {/* Toast Container */}
      <Box 
        position="fixed" 
        bottom="4" 
        right="4" 
        zIndex="50"
        display="flex"
        flexDirection="column"
        gap="2"
      >
        {toasts.map((t) => (
          <Box 
            key={t.id}
            padding="4"
            borderRadius="md"
            background="surface-container-highest"
            border={`1px solid \${t.type === 'error' ? 'var(--error)' : 'var(--accent)'}`}
            boxShadow="lg"
            minWidth="300px"
            display="flex"
            alignItems="center"
            justifyContent="space-between"
            className="animate-in fade-in slide-in-from-right-4"
          >
            <Stack direction="row" gap="3" alignItems="center">
              {t.type === 'success' ? <CheckCircle size={18} color="green" /> : <AlertCircle size={18} color={t.type === 'error' ? 'red' : 'cyan'} />}
              <Text weight="medium">{t.message}</Text>
            </Stack>
            <Box 
              onClick={() => setToasts((prev) => prev.filter((toast) => toast.id !== t.id))}
              cursor="pointer"
              color="muted"
            >
              <X size={16} />
            </Box>
          </Box>
        ))}
      </Box>
    </ToastContext.Provider>
  );
}

export function useToast() {
  const context = useContext(ToastContext);
  if (context === undefined) {
    throw new Error("useToast must be used within a ToastProvider");
  }
  return context;
}
