<template>
  <footer class="footer">
    <div class="footer-inner">
      <p>&copy; {{ year }} Kuonji</p>
    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>

<style scoped>
.footer {
  border-top: 1px solid var(--color-border);
  margin-top: var(--space-20);
  padding: var(--space-8) 0;
}

.footer-inner {
  max-width: var(--max-width-wide);
  margin: 0 auto;
  padding: 0 var(--space-6);
  display: flex;
  justify-content: center;
  align-items: center;
}

.footer p {
  font-size: var(--text-sm);
  color: var(--color-muted);
  margin: 0;
}
</style>
