import { argon2id } from "hash-wasm";

const ARGON2_PARAMS = {
  iterations: 1,
  memory: 64 * 1024, // 64MB
  parallelism: 4,
  hashLength: 32,
};

/**
 * Derives a key from a password and salt using Argon2id.
 */
export async function deriveKey(password: string, salt: Uint8Array): Promise<Uint8Array> {
  return argon2id({
    password,
    salt,
    iterations: ARGON2_PARAMS.iterations,
    memorySize: ARGON2_PARAMS.memory,
    parallelism: ARGON2_PARAMS.parallelism,
    hashLength: ARGON2_PARAMS.hashLength,
    outputType: "binary",
  });
}

/**
 * Encrypts data using AES-GCM.
 * Returns Uint8Array containing [nonce (12 bytes) + ciphertext + tag (16 bytes)].
 */
export async function encrypt(key: Uint8Array, plaintext: Uint8Array): Promise<Uint8Array> {
  const cryptoKey = await window.crypto.subtle.importKey(
    "raw",
    new Uint8Array(key),
    { name: "AES-GCM" },
    false,
    ["encrypt"]
  );

  const nonce = window.crypto.getRandomValues(new Uint8Array(12));
  const encrypted = await window.crypto.subtle.encrypt(
    { name: "AES-GCM", iv: nonce },
    cryptoKey,
    new Uint8Array(plaintext)
  );

  const result = new Uint8Array(nonce.length + encrypted.byteLength);
  result.set(nonce);
  result.set(new Uint8Array(encrypted), nonce.length);
  return result;
}

/**
 * Decrypts data using AES-GCM.
 */
export async function decrypt(key: Uint8Array, ciphertextWithNonce: Uint8Array): Promise<Uint8Array> {
  const cryptoKey = await window.crypto.subtle.importKey(
    "raw",
    new Uint8Array(key),
    { name: "AES-GCM" },
    false,
    ["decrypt"]
  );

  const nonce = ciphertextWithNonce.slice(0, 12);
  const data = ciphertextWithNonce.slice(12);

  const decrypted = await window.crypto.subtle.decrypt(
    { name: "AES-GCM", iv: nonce },
    cryptoKey,
    data
  );

  return new Uint8Array(decrypted);
}

/**
 * Helper to convert hex string to Uint8Array.
 */
export function hexToBytes(hex: string): Uint8Array {
  const bytes = new Uint8Array(hex.length / 2);
  for (let i = 0; i < bytes.length; i++) {
    bytes[i] = parseInt(hex.substr(i * 2, 2), 16);
  }
  return bytes;
}

/**
 * Helper to convert Uint8Array to hex string.
 */
export function bytesToHex(bytes: Uint8Array): string {
  return Array.from(bytes)
    .map((b) => b.toString(16).padStart(2, "0"))
    .join("");
}
