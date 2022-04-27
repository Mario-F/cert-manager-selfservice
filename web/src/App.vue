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
  CHeaderToggler,
  CCollapse,
  CHeaderNav,
  CNavLink,
  CDropdown,
  CDropdownItem,
} from '@coreui/vue'
import { CIcon } from '@coreui/icons-vue'
import * as icons from '@coreui/icons/js/free'

let visible = true
</script>

<template>
  <div>
    <CSidebar visible>
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

    <div class="wrapper d-flex flex-column min-vh-100 bg-light">
      <CHeader>
        <CContainer fluid>
          <CHeaderBrand href="#">Header</CHeaderBrand>
          <CHeaderToggler @click="visible = !visible" />
          <CCollapse class="header-collapse" :visible="visible">
            <CHeaderNav>
              <CNavItem>
                <CNavLink href="#" active> Home </CNavLink>
              </CNavItem>
              <CNavItem>
                <CNavLink href="#">Link</CNavLink>
              </CNavItem>
              <CDropdown togglerText="Dropdown" variant="nav-item">
                <CDropdownItem href="#">Action</CDropdownItem>
                <CDropdownItem href="#">Another action</CDropdownItem>
                <CDropdownItem href="#">Something else here</CDropdownItem>
              </CDropdown>
              <CNavItem>
                <CNavLink href="#" disabled> Disabled </CNavLink>
              </CNavItem>
            </CHeaderNav>
            <!-- <form class="d-flex">
              <input class="form-control me-2" type="search" placeholder="Search">
              <button class="btn btn-outline-success" type="submit">Search</button>
            </form> -->
          </CCollapse>
        </CContainer>
      </CHeader>
      <div class="body flex-grow-1 px-3">
        <CContainer lg>
          <RouterView />
        </CContainer>
      </div>
    </div>
  </div>
</template>

<style></style>
