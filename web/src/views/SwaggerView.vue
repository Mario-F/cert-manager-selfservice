<script setup lang="ts">
import { onMounted } from 'vue'
import 'swagger-ui/dist/swagger-ui.css'
import SwaggerUI from 'swagger-ui'

onMounted(() => {
  SwaggerUI({
    dom_id: '#swagger-ui',
    url: '/api/spec/v1',
  })
})
</script>

<template>
  <div id="swagger-ui"></div>
</template>
