<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { CSidebar, CSidebarBrand, CSidebarNav, CNavItem, CBadge, CNavGroup, CSidebarToggler } from '@coreui/vue'
</script>

<template>
  <div>
    <CSidebar visible>
      <CSidebarBrand>Sidebar Brand</CSidebarBrand>
      <CSidebarNav>
        <li class="nav-title">Nav Title</li>
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" icon="cil-speedometer" />
          Nav item
        </CNavItem>
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" icon="cil-speedometer" />
          With badge
          <CBadge color="primary ms-auto">NEW</CBadge>
        </CNavItem>
        <CNavGroup>
          <template #togglerContent> <CIcon customClassName="nav-icon" icon="cil-puzzle" /> Nav dropdown </template>
          <CNavItem href="#"> <CIcon customClassName="nav-icon" icon="cil-puzzle" /> Nav dropdown item </CNavItem>
          <CNavItem href="#"> <CIcon customClassName="nav-icon" icon="cil-puzzle" /> Nav dropdown item </CNavItem>
        </CNavGroup>
      </CSidebarNav>
      <CSidebarToggler />
    </CSidebar>

    <nav>
      <RouterLink to="/">Home</RouterLink>
      <RouterLink to="/about">About</RouterLink>
    </nav>
  </div>

  <RouterView />
</template>

<style></style>
