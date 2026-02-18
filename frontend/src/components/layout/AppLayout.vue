<template>
  <div class="app-layout">
    <div class="app-layout__list">
      <slot name="list" />
    </div>
    <div class="app-layout__form">
      <slot name="form" />
    </div>
  </div>
</template>
