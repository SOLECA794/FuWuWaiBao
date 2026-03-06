# Vue 3 + Vite

This template should help get you started developing with Vue 3 in Vite. The template uses Vue 3 `<script setup>` SFCs, check out the [script setup docs](https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup) to learn more.

Learn more about IDE Support for Vue in the [Vue Docs Scaling up Guide](https://vuejs.org/guide/scaling-up/tooling.html#ide-support).

## Environment

Copy `.env.example` to `.env.local` and adjust backend address if needed:

```bash
cp .env.example .env.local
```

`VITE_API_BASE` defaults to `http://localhost:18080`.

## Project Structure

- `src/config/api.js`: backend base URL config.
- `src/services/teacherApi.js`: teacher-side API request layer.
- `src/styles/theme.css`: global theme tokens and base styles.
