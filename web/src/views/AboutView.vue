<script setup lang="ts">
import { CContainer, CRow, CCol, CCard, CCardBody, CCardTitle, CCardText, CHeaderBrand } from '@coreui/vue'
import { CIcon } from '@coreui/icons-vue'
import * as icons from '@coreui/icons'
import { onMounted, ref } from 'vue'
import { DefaultService } from '../gen/api'

const info = {
  version: ref('...'),
}

onMounted(async () => {
  const resInfo = await DefaultService.getInfo()
  info.version.value = resInfo.version
})
</script>

<template>
  <CContainer fluid>
    <CRow>
      <CHeaderBrand>About</CHeaderBrand>
    </CRow>
    <CRow>
      <CCol>
        <CCard>
          <CCardBody>
            <CCardTitle>Version - {{ info.version.value }}</CCardTitle>
            <CCardText>Check out the repository at Github.</CCardText>
            <a
              class="btn btn-outline-primary"
              href="https://github.com/Mario-F/cert-manager-selfservice"
              target="_blank"
              role="button"
              ><CIcon :icon="icons.cibGithub" /> Github</a
            >
          </CCardBody>
        </CCard>
      </CCol>
    </CRow>
  </CContainer>
</template>
