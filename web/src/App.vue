<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
import {
  CSidebar,
  CSidebarBrand,
  CSidebarNav,
  CNavItem,
  CBadge,
  CNavGroup,
  CSidebarToggler,
  CContainer,
  CHeader,
  CHeaderBrand,
  CHeaderNav,
  CDropdownMenu,
  CDropdownToggle,
  CNavLink,
  CDropdown,
  CDropdownItem,
} from '@coreui/vue'
import { CIcon } from '@coreui/icons-vue'
import * as icons from '@coreui/icons/js/free'
import { ref } from 'vue'

const sidebar = {
  visible: ref(true),
  unfoldable: ref(false),
}
</script>

<template>
  <CSidebar position="fixed" :unfoldable="sidebar.unfoldable.value" :visible="sidebar.visible.value">
    <CSidebarBrand>Sidebar Brand</CSidebarBrand>
    <CSidebarNav>
      <li class="nav-title">Nav Title</li>
      <RouterLink to="/">
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilSpeedometer" />
          Home
        </CNavItem>
      </RouterLink>
      <RouterLink to="/about">
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilSpeedometer" />
          About
          <CBadge color="primary" class="ms-auto">NEW</CBadge>
        </CNavItem>
      </RouterLink>
      <RouterLink to="/swagger">
        <CNavItem href="#">
          <CIcon customClassName="nav-icon" :icon="icons.cilSpeedometer" />
          API
          <CBadge color="primary" class="ms-auto">Swagger</CBadge>
        </CNavItem>
      </RouterLink>
      <CNavGroup>
        <template #togglerContent> <CIcon customClassName="nav-icon" :icon="icons.cilPuzzle" /> Nav dropdown </template>
        <CNavItem href="#"> <CIcon customClassName="nav-icon" :icon="icons.cilPuzzle" /> Nav dropdown item </CNavItem>
        <CNavItem href="#"> <CIcon customClassName="nav-icon" :icon="icons.cilPuzzle" /> Nav dropdown item </CNavItem>
      </CNavGroup>
    </CSidebarNav>
    <CSidebarToggler class="d-none d-lg-flex" @click="sidebar.unfoldable.value = !sidebar.unfoldable.value" />
  </CSidebar>

  <div class="wrapper d-flex flex-column min-vh-100 bg-light">
    <CHeader>
      <CContainer fluid>
        <CHeaderBrand href="#">Header</CHeaderBrand>
        <CHeaderNav variant="pills">
          <CNavItem>
            <RouterLink to="/"><CNavLink>Home</CNavLink></RouterLink>
          </CNavItem>
          <CNavItem>
            <RouterLink to="/about"><CNavLink>About</CNavLink></RouterLink>
          </CNavItem>
          <CDropdown togglerText="Dropdown" variant="nav-item">
            <CDropdownToggle color="primary">Dropdown Button</CDropdownToggle>
            <CDropdownMenu>
              <CDropdownItem href="#">Action</CDropdownItem>
              <CDropdownItem href="#">Another action</CDropdownItem>
              <CDropdownItem href="#">Something else here</CDropdownItem>
            </CDropdownMenu>
          </CDropdown>
          <CNavItem>
            <CNavLink href="#" disabled> Disabled </CNavLink>
          </CNavItem>
        </CHeaderNav>
      </CContainer>
    </CHeader>
    <div class="body flex-grow-1 px-3">
      <CContainer lg>
        <RouterView />
      </CContainer>
    </div>
  </div>
</template>

<style lang="scss">
@import 'styles/style';
</style>
