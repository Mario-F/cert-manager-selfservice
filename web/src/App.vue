<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import { CSidebar, CSidebarBrand, CSidebarNav, CNavItem, CBadge, CNavGroup, CSidebarToggler } from '@coreui/vue'
import { CIcon } from '@coreui/icons-vue'
import * as icons from '@coreui/icons/js/free'
</script>

<template>
  <div>
    <CSidebar visible>
      <CSidebarBrand>Sidebar Brand</CSidebarBrand>
      <CSidebarNav>
        <li class="nav-title">Nav Title</li>
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilSpeedometer" />
          Nav item
        </CNavItem>
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilSpeedometer" />
          With badge
          <CBadge color="primary" class="ms-auto">NEW</CBadge>
        </CNavItem>
        <CNavGroup>
          <template #togglerContent>
            <CIcon customClassName="nav-icon" :icon="icons.cilPuzzle" /> Nav dropdown
          </template>
          <CNavItem href="#"> <CIcon customClassName="nav-icon" :icon="icons.cilPuzzle" /> Nav dropdown item </CNavItem>
          <CNavItem href="#"> <CIcon customClassName="nav-icon" :icon="icons.cilPuzzle" /> Nav dropdown item </CNavItem>
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
