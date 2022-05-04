<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { CSidebar, CSidebarBrand, CSidebarNav, CNavItem, CBadge, CSidebarToggler } from '@coreui/vue'
import { CIcon } from '@coreui/icons-vue'
import * as icons from '@coreui/icons'
import { ref } from 'vue'

const sidebar = {
  unfoldable: ref(false),
}
</script>

<template>
  <CSidebar position="fixed" :unfoldable="sidebar.unfoldable.value" visible>
    <CSidebarBrand>Cert Manager Selfservice</CSidebarBrand>
    <CSidebarNav>
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

  <div class="router-wrapper">
    <RouterView />
  </div>
</template>

<style lang="scss">
@import 'styles/style';
</style>
