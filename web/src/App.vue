<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { CSidebar, CSidebarBrand, CSidebarNav, CNavItem, CBadge, CSidebarToggler } from '@coreui/vue'
import { CIcon } from '@coreui/icons-vue'
import * as icons from '@coreui/icons'
import { ref, onMounted } from 'vue'
import { DefaultApi } from './gen/api'

const sidebar = {
  unfoldable: ref(false),
}

const info = {
  version: ref('...'),
}

onMounted(async () => {
  const api = new DefaultApi()
  const resInfo = await api.infoGet()
  info.version.value = resInfo.data.version
})
</script>

<template>
  <CSidebar position="fixed" :unfoldable="sidebar.unfoldable.value" visible>
    <CSidebarBrand>Cert Manager Selfservice</CSidebarBrand>
    <CSidebarNav>
      <li class="nav-title">{{ info.version.value }}</li>
      <RouterLink to="/">
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilSpeedometer" />
          Dashboard
        </CNavItem>
      </RouterLink>
      <RouterLink to="/swagger">
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cibSwagger" />
          API
          <CBadge color="primary" class="ms-auto">Swagger</CBadge>
        </CNavItem>
      </RouterLink>
      <RouterLink to="/about">
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilInfo" />
          About
        </CNavItem>
      </RouterLink>
    </CSidebarNav>
    <CSidebarToggler class="d-none d-lg-flex" @click="sidebar.unfoldable.value = !sidebar.unfoldable.value" />
  </CSidebar>

  <div class="wrapper d-flex flex-column min-vh-100 bg-light">
    <RouterView />
  </div>
</template>

<style lang="scss">
@import 'styles/style';
</style>
