<template>
  <div class="browser-view" ref="rootEl"
       @dragover.prevent="onDragOver"
       @dragleave.self="onDragLeave"
       @drop.prevent="onDrop">

    <!-- ── Header ──────────────────────────────────────────────── -->
    <div class="browser-hd">
      <div class="browser-hd__left">
        <div class="browser-prov-icon" :class="`browser-prov-icon--${conn.provider}`">
          <ProviderIcon :provider="conn.provider" :size="16" />
        </div>
        <div style="min-width:0">
          <div class="browser-conn-name">
            {{ conn.name }}
            <BaseBadge :provider="conn.provider" />
          </div>
          <div class="browser-conn-bucket">{{ conn.bucket }}</div>
          <div class="breadcrumbs" v-if="currentPrefix">
            <button class="bread-item" @click="navigateTo('')">root</button>
            <template v-for="(crumb, i) in breadcrumbs" :key="i">
              <span class="bread-sep">/</span>
              <button
                class="bread-item"
                :style="i === breadcrumbs.length - 1 ? 'color:var(--text-2);cursor:default' : ''"
                @click="i < breadcrumbs.length - 1 && navigateTo(crumb.prefix)"
              >{{ crumb.label }}</button>
            </template>
          </div>
        </div>
      </div>

      <div class="browser-hd__actions">
        <button class="icon-btn" :style="showStats ? 'background:var(--accent-bg);color:var(--accent);border-color:var(--accent-ring)' : ''" @click="toggleStats" title="Bucket stats">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="20" x2="18" y2="10"/><line x1="12" y1="20" x2="12" y2="4"/><line x1="6" y1="20" x2="6" y2="14"/>
          </svg>
        </button>
        <button class="icon-btn" :disabled="loading" @click="refresh" title="Refresh (R)">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" :style="loading ? 'animation:spin .6s linear infinite' : ''">
            <polyline points="23 4 23 10 17 10"/><polyline points="1 20 1 14 7 14"/>
            <path d="M3.51 9a9 9 0 0 1 14.85-3.36L23 10M1 14l4.64 4.36A9 9 0 0 0 20.49 15"/>
          </svg>
        </button>
        <button class="icon-btn danger" @click="$emit('delete')" title="Delete connection">
          <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
            <path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4h6v2"/>
          </svg>
        </button>
      </div>
    </div>

    <!-- ── Stats bar ────────────────────────────────────────────── -->
    <transition name="slide-down">
      <div v-if="showStats" class="stats-bar">
        <template v-if="statsLoading">
          <div class="stat-item">
            <div class="skeleton-item" style="width:60px;height:22px;border-radius:4px"></div>
            <span class="stat-lbl">loading…</span>
          </div>
        </template>
        <template v-else-if="stats">
          <div class="stat-item">
            <span class="stat-val">{{ stats.object_count.toLocaleString() }}</span>
            <span class="stat-lbl">{{ stats.truncated ? 'objects (est.)' : 'objects' }}</span>
          </div>
          <div class="stat-item">
            <span class="stat-val">{{ formatSize(stats.total_size) }}</span>
            <span class="stat-lbl">total size</span>
          </div>
        </template>
        <template v-else-if="statsError">
          <span style="font-size:12px;color:var(--muted)">{{ statsError }}</span>
        </template>
      </div>
    </transition>

    <!-- ── Toolbar ──────────────────────────────────────────────── -->
    <div class="browser-toolbar">
      <div class="search-field">
        <span class="search-field__icon">
          <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="11" cy="11" r="8"/><line x1="21" y1="21" x2="16.65" y2="16.65"/>
          </svg>
        </span>
        <input ref="searchInput" class="search-field__input" v-model="searchQuery" placeholder="Search files… (/)" @keydown.escape.stop="searchQuery = ''" />
        <button v-if="searchQuery" class="search-field__clear" @click="searchQuery = ''">×</button>
      </div>

      <div class="toolbar-spacer"></div>

      <!-- New folder -->
      <button class="icon-btn" @click="showFolderModal = true" title="New folder">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
          <line x1="12" y1="11" x2="12" y2="17"/><line x1="9" y1="14" x2="15" y2="14"/>
        </svg>
      </button>

      <!-- Upload -->
      <label class="upload-label" title="Upload files">
        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/>
          <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/>
        </svg>
        Upload
        <input type="file" multiple style="display:none" @change="onFileInput" />
      </label>
    </div>

    <!-- ── Upload progress ──────────────────────────────────────── -->
    <transition name="slide-down">
      <div v-if="uploading" class="upload-progress">
        <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" style="animation:spin .6s linear infinite" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="2" x2="12" y2="6"/><line x1="12" y1="18" x2="12" y2="22"/>
          <line x1="4.93" y1="4.93" x2="7.76" y2="7.76"/><line x1="16.24" y1="16.24" x2="19.07" y2="19.07"/>
          <line x1="2" y1="12" x2="6" y2="12"/><line x1="18" y1="12" x2="22" y2="12"/>
          <line x1="4.93" y1="19.07" x2="7.76" y2="16.24"/><line x1="16.24" y1="7.76" x2="19.07" y2="4.93"/>
        </svg>
        <span>Uploading {{ uploadingCount }} file{{ uploadingCount !== 1 ? 's' : '' }}…</span>
        <div class="progress-bar"><div class="progress-fill" style="width:100%"></div></div>
      </div>
    </transition>

    <!-- ── Selection action bar ─────────────────────────────────── -->
    <transition name="slide-down">
      <div v-if="selected.size > 0" class="selection-bar">
        <span class="selection-bar__count">{{ selected.size }} selected</span>
        <button class="base-btn base-btn--ghost" style="font-size:12px;padding:5px 10px" @click="bulkDownload" :disabled="bulkWorking">
          Download all
        </button>
        <button class="base-btn base-btn--danger" style="font-size:12px;padding:5px 10px;border:1px solid var(--danger)" @click="bulkDelete" :disabled="bulkWorking">
          Delete all
        </button>
        <button class="icon-btn" style="width:26px;height:26px;margin-left:auto" title="Deselect all" @click="selected.clear(); selected = new Set(selected)">
          <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
      </div>
    </transition>

    <!-- ── Body ─────────────────────────────────────────────────── -->
    <div class="browser-body" ref="bodyEl">

      <!-- Drag overlay -->
      <div v-if="isDragging" class="drop-overlay">
        <div class="drop-overlay__inner">
          <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" style="margin-bottom:8px">
            <polyline points="16 16 12 12 8 16"/><line x1="12" y1="12" x2="12" y2="21"/>
            <path d="M20.39 18.39A5 5 0 0 0 18 9h-1.26A8 8 0 1 0 3 16.3"/>
          </svg>
          Drop to upload here
        </div>
      </div>

      <!-- Loading skeleton -->
      <table v-if="loading && entries.length === 0" class="file-table">
        <thead><tr><th class="col-check"></th><th>Name</th><th>Size</th><th>Modified</th><th></th></tr></thead>
        <tbody>
          <tr v-for="i in 10" :key="i">
            <td class="col-check"></td>
            <td colspan="4"><div class="skeleton-item" :style="`height:15px;border-radius:3px;width:${30 + (i * 37 % 55)}%`"></div></td>
          </tr>
        </tbody>
      </table>

      <!-- Browse error -->
      <div v-else-if="browseError && entries.length === 0" class="empty-state">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.4;margin-bottom:6px">
          <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
        </svg>
        <p style="font-size:13px;color:var(--text-2);max-width:320px">{{ browseError }}</p>
        <button class="base-btn base-btn--ghost" @click="refresh" style="font-size:12px;padding:6px 12px;margin-top:4px">Retry</button>
      </div>

      <!-- Empty -->
      <div v-else-if="!loading && filteredEntries.length === 0 && !nextPageToken" class="empty-state">
        <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.4;margin-bottom:6px">
          <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z"/>
        </svg>
        <p>{{ searchQuery ? 'No files match your search.' : 'This folder is empty.' }}</p>
        <button v-if="searchQuery" class="base-btn base-btn--ghost" @click="searchQuery = ''" style="font-size:12px;padding:6px 12px;margin-top:4px">Clear search</button>
      </div>

      <!-- File table -->
      <table v-else class="file-table">
        <thead>
          <tr>
            <th class="col-check">
              <input
                type="checkbox"
                :checked="allFilesSelected"
                :indeterminate="someFilesSelected && !allFilesSelected"
                @change="toggleSelectAll"
                title="Select all files"
              />
            </th>
            <th class="sortable" @click="cycleSort('name')">
              Name <span class="sort-indicator" :class="{active: sortKey==='name'}">{{ sortKey==='name' ? (sortDir==='asc'?'↑':'↓') : '↕' }}</span>
            </th>
            <th class="sortable" @click="cycleSort('size')">
              Size <span class="sort-indicator" :class="{active: sortKey==='size'}">{{ sortKey==='size' ? (sortDir==='asc'?'↑':'↓') : '↕' }}</span>
            </th>
            <th class="sortable" @click="cycleSort('date')">
              Modified <span class="sort-indicator" :class="{active: sortKey==='date'}">{{ sortKey==='date' ? (sortDir==='asc'?'↑':'↓') : '↕' }}</span>
            </th>
            <th style="width:132px"></th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="entry in filteredEntries"
            :key="entry.name"
            :class="{ 'is-dir': entry.type === 'dir', 'is-selected': selected.has(entry.name) }"
          >
            <td class="col-check">
              <input v-if="entry.type === 'file'" type="checkbox" :checked="selected.has(entry.name)" @change="toggleSelect(entry.name)" />
            </td>
            <td>
              <div class="file-name" :style="entry.type === 'dir' ? 'cursor:pointer' : ''" @click="entry.type === 'dir' && navigateTo(entry.name)">
                <svg v-if="entry.type === 'dir'" class="file-icon" width="13" height="13" viewBox="0 0 24 24" fill="currentColor" stroke="none" style="color:var(--aws)">
                  <path d="M22 19a2 2 0 0 1-2 2H4a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h5l2 3h9a2 2 0 0 1 2 2z" opacity=".8"/>
                </svg>
                <svg v-else class="file-icon" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="13 2 13 9 20 9"/>
                </svg>
                {{ entry.display }}
              </div>
            </td>
            <td class="file-size">{{ entry.type === 'dir' ? '—' : formatSize(entry.size) }}</td>
            <td class="file-date">{{ entry.type === 'dir' ? '—' : formatDate(entry.updated) }}</td>
            <td class="file-actions">
              <!-- Copy path -->
              <button class="row-btn" @click.stop="copyPath(entry)" title="Copy path">
                <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/><path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"/>
                </svg>
              </button>
              <template v-if="entry.type === 'file'">
                <!-- Rename -->
                <button class="row-btn" @click.stop="openRename(entry)" title="Rename / Move">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                    <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                  </svg>
                </button>
                <!-- Download -->
                <button class="row-btn" @click.stop="download(entry)" title="Download">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/>
                  </svg>
                </button>
                <!-- Preview / info -->
                <button class="row-btn" @click.stop="openPreview(entry)" title="Preview">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/><circle cx="12" cy="12" r="3"/>
                  </svg>
                </button>
                <!-- Metadata -->
                <button class="row-btn" @click.stop="openMeta(entry)" title="Metadata">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/>
                  </svg>
                </button>
                <!-- Delete -->
                <button class="row-btn danger" @click.stop="confirmDelete(entry)" title="Delete">
                  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <polyline points="3 6 5 6 21 6"/><path d="M19 6l-1 14a2 2 0 0 1-2 2H8a2 2 0 0 1-2-2L5 6"/>
                    <path d="M10 11v6"/><path d="M14 11v6"/><path d="M9 6V4h6v2"/>
                  </svg>
                </button>
              </template>
            </td>
          </tr>

          <!-- Infinite scroll sentinel -->
          <tr v-if="nextPageToken" ref="sentinel">
            <td colspan="5" style="padding:12px 22px;text-align:center;font-size:12px;color:var(--muted)">
              <span v-if="loadingMore">Loading more…</span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ── Preview panel ────────────────────────────────────────── -->
    <transition name="slide-right">
      <div v-if="previewEntry && !metaEntry" class="preview-panel">
        <div class="preview-hd">
          <span class="preview-hd__name">{{ previewEntry.display }}</span>
          <button class="preview-close" @click="closePreview">×</button>
        </div>
        <div class="preview-body">
          <div v-if="previewLoading" class="preview-unsupported">
            <div class="base-btn__spinner" style="width:20px;height:20px;border-width:2px"></div>
          </div>
          <img v-else-if="isImage(previewEntry) && previewUrl" :src="previewUrl" class="preview-img" @error="previewLoadError=true" />
          <pre v-else-if="isText(previewEntry) && previewContent" class="preview-text">{{ previewContent }}</pre>
          <div v-else class="preview-unsupported">
            <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round" style="opacity:.3">
              <path d="M13 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V9z"/><polyline points="13 2 13 9 20 9"/>
            </svg>
            <p>{{ previewLoadError ? 'Failed to load preview.' : 'No preview for this file type.' }}</p>
          </div>
        </div>
        <div class="preview-ft">
          <span class="preview-meta">
            {{ formatSize(previewEntry.size) }}
            <template v-if="previewEntry.content_type"> · {{ previewEntry.content_type }}</template>
          </span>
          <button class="base-btn base-btn--ghost" @click="download(previewEntry)" style="font-size:12px;padding:5px 10px">Download</button>
        </div>
      </div>
    </transition>

    <!-- ── Metadata panel ───────────────────────────────────────── -->
    <transition name="slide-right">
      <div v-if="metaEntry" class="preview-panel">
        <div class="preview-hd">
          <span class="preview-hd__name">{{ metaEntry.display }} — Metadata</span>
          <button class="preview-close" @click="metaEntry = null">×</button>
        </div>
        <div class="preview-body" style="padding:0">
          <div v-if="metaLoading" class="preview-unsupported">
            <div class="base-btn__spinner" style="width:20px;height:20px;border-width:2px"></div>
          </div>
          <div v-else-if="metaData" style="padding:16px;display:flex;flex-direction:column;gap:14px">
            <!-- Content-Type -->
            <div class="meta-field">
              <label class="meta-label">Content-Type</label>
              <input class="base-input" v-model="metaEdit.content_type" style="font-size:12px;padding:6px 10px" />
            </div>
            <!-- Cache-Control -->
            <div class="meta-field">
              <label class="meta-label">Cache-Control</label>
              <input class="base-input" v-model="metaEdit.cache_control" style="font-size:12px;padding:6px 10px" />
            </div>
            <!-- Custom metadata -->
            <div class="meta-field">
              <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:6px">
                <label class="meta-label" style="margin-bottom:0">Custom Metadata</label>
                <button class="base-btn base-btn--ghost" @click="addMetaRow" style="font-size:11px;padding:3px 8px">+ Add</button>
              </div>
              <div v-for="(pair, i) in metaRows" :key="i" class="meta-row">
                <input class="base-input" v-model="pair.key" placeholder="key" style="font-size:11px;padding:5px 8px;flex:1" />
                <input class="base-input" v-model="pair.val" placeholder="value" style="font-size:11px;padding:5px 8px;flex:2" />
                <button class="row-btn danger" @click="metaRows.splice(i,1)" style="opacity:1;flex-shrink:0">
                  <svg width="11" height="11" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                    <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
                  </svg>
                </button>
              </div>
              <p v-if="metaRows.length === 0" style="font-size:11px;color:var(--muted)">No custom metadata.</p>
            </div>
            <!-- Read-only info -->
            <div style="padding:10px;background:var(--surface-2);border-radius:var(--r-sm);font-size:11px;color:var(--muted);line-height:1.8">
              <div>Size: <strong style="color:var(--text-2)">{{ formatSize(metaData.size) }}</strong></div>
              <div v-if="metaData.etag">ETag: <strong style="color:var(--text-2);font-family:var(--mono)">{{ metaData.etag }}</strong></div>
              <div v-if="metaData.md5">MD5: <strong style="color:var(--text-2);font-family:var(--mono)">{{ metaData.md5 }}</strong></div>
              <div v-if="metaData.updated">Modified: <strong style="color:var(--text-2)">{{ formatDate(metaData.updated) }}</strong></div>
            </div>
          </div>
          <div v-else-if="metaError" class="preview-unsupported" style="font-size:12px">{{ metaError }}</div>
        </div>
        <div class="preview-ft">
          <button class="base-btn base-btn--ghost" @click="metaEntry = null" style="font-size:12px;padding:5px 10px">Cancel</button>
          <button class="base-btn base-btn--primary" @click="saveMeta" :disabled="metaSaving" style="font-size:12px;padding:5px 12px">
            {{ metaSaving ? 'Saving…' : 'Save' }}
          </button>
        </div>
      </div>
    </transition>

    <!-- ── Modals ────────────────────────────────────────────────── -->

    <!-- New folder -->
    <BaseModal :open="showFolderModal" title="New Folder" @update:open="showFolderModal = false">
      <div style="display:flex;flex-direction:column;gap:10px">
        <label class="form-label">Folder name</label>
        <input
          ref="folderInput"
          class="base-input"
          v-model="newFolderName"
          placeholder="e.g. images"
          @keydown.enter="createFolder"
          @keydown.escape.stop="showFolderModal = false"
        />
        <p class="form-hint">A placeholder file (<code style="font-family:var(--mono)">.keep</code>) will be uploaded inside the new folder.</p>
      </div>
      <template #footer>
        <button class="base-btn base-btn--ghost" @click="showFolderModal = false">Cancel</button>
        <button class="base-btn base-btn--primary" @click="createFolder" :disabled="!newFolderName.trim()">Create</button>
      </template>
    </BaseModal>

    <!-- Rename / Move -->
    <BaseModal :open="showRenameModal" :title="`Rename: ${renameEntry?.display ?? ''}`" @update:open="showRenameModal = false">
      <div style="display:flex;flex-direction:column;gap:10px">
        <label class="form-label">New name (within current folder)</label>
        <input
          ref="renameInput"
          class="base-input"
          v-model="renameTarget"
          @keydown.enter="doRename"
          @keydown.escape.stop="showRenameModal = false"
        />
        <p class="form-hint">Current: <code style="font-family:var(--mono)">{{ renameEntry?.name }}</code></p>
      </div>
      <template #footer>
        <button class="base-btn base-btn--ghost" @click="showRenameModal = false">Cancel</button>
        <button class="base-btn base-btn--primary" @click="doRename" :disabled="!renameTarget.trim() || renaming">
          {{ renaming ? 'Moving…' : 'Move' }}
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'
import BaseBadge    from '../ui/BaseBadge.vue'
import BaseModal    from '../ui/BaseModal.vue'
import ProviderIcon from '../ui/ProviderIcon.vue'
import { useConnections }   from '../../composables/useConnections.js'
import { useToast }         from '../../composables/useToast.js'
import { useConfirm }       from '../../composables/useConfirm.js'

const props = defineProps({ conn: { type: Object, required: true } })
defineEmits(['delete'])

const { browseObjects, getDownloadURL, deleteObject, copyObject, uploadObjects, getBucketStats, getObjectMetadata, updateObjectMetadata } = useConnections()
const toast   = useToast()
const confirm = useConfirm()

// ── Core state ──────────────────────────────────────────────────
const currentPrefix  = ref('')
const entries        = ref([])
const nextPageToken  = ref('')
const loading        = ref(false)
const loadingMore    = ref(false)
const browseError    = ref('')

const searchQuery    = ref('')
const sortKey        = ref('name')
const sortDir        = ref('asc')

// ── Stats ───────────────────────────────────────────────────────
const stats        = ref(null)
const statsLoading = ref(false)
const statsError   = ref('')
const statsLoaded  = ref(false)
const showStats    = ref(false)

// ── Upload ──────────────────────────────────────────────────────
const uploading      = ref(false)
const uploadingCount = ref(0)
const isDragging     = ref(false)

// ── Bulk select ─────────────────────────────────────────────────
let selected = ref(new Set())
const bulkWorking = ref(false)

const fileEntries = computed(() => entries.value.filter(e => e.type === 'file'))
const allFilesSelected  = computed(() => fileEntries.value.length > 0 && fileEntries.value.every(e => selected.value.has(e.name)))
const someFilesSelected = computed(() => fileEntries.value.some(e => selected.value.has(e.name)))

function toggleSelect(name) {
  const s = new Set(selected.value)
  s.has(name) ? s.delete(name) : s.add(name)
  selected.value = s
}

function toggleSelectAll() {
  if (allFilesSelected.value) {
    selected.value = new Set()
  } else {
    selected.value = new Set(fileEntries.value.map(e => e.name))
  }
}

// ── Preview ─────────────────────────────────────────────────────
const previewEntry    = ref(null)
const previewUrl      = ref('')
const previewContent  = ref('')
const previewLoading  = ref(false)
const previewLoadError = ref(false)

// ── Metadata editor ─────────────────────────────────────────────
const metaEntry   = ref(null)
const metaData    = ref(null)
const metaEdit    = ref({ content_type: '', cache_control: '' })
const metaRows    = ref([]) // [{ key, val }]
const metaLoading = ref(false)
const metaError   = ref('')
const metaSaving  = ref(false)

// ── Modals ──────────────────────────────────────────────────────
const showFolderModal = ref(false)
const newFolderName   = ref('')
const folderInput     = ref(null)

const showRenameModal = ref(false)
const renameEntry     = ref(null)
const renameTarget    = ref('')
const renameInput     = ref(null)
const renaming        = ref(false)

// ── DOM refs ────────────────────────────────────────────────────
const searchInput = ref(null)
const bodyEl      = ref(null)
const sentinel    = ref(null)
let   observer    = null

// ── Breadcrumbs ─────────────────────────────────────────────────
const breadcrumbs = computed(() => {
  if (!currentPrefix.value) return []
  const parts = currentPrefix.value.replace(/\/$/, '').split('/')
  return parts.map((label, i) => ({
    label,
    prefix: parts.slice(0, i + 1).join('/') + '/',
  }))
})

// ── Filtered + sorted entries ───────────────────────────────────
const filteredEntries = computed(() => {
  let list = entries.value
  if (searchQuery.value.trim()) {
    const q = searchQuery.value.trim().toLowerCase()
    list = list.filter(e => e.display.toLowerCase().includes(q))
  }
  const dirs  = list.filter(e => e.type === 'dir')
  const files = list.filter(e => e.type === 'file')
  const sortFn = (a, b) => {
    let va, vb
    if (sortKey.value === 'size')      { va = a.size ?? 0; vb = b.size ?? 0 }
    else if (sortKey.value === 'date') { va = a.updated ? new Date(a.updated).getTime() : 0; vb = b.updated ? new Date(b.updated).getTime() : 0 }
    else                               { va = a.display.toLowerCase(); vb = b.display.toLowerCase() }
    if (va < vb) return sortDir.value === 'asc' ? -1 : 1
    if (va > vb) return sortDir.value === 'asc' ?  1 : -1
    return 0
  }
  return [...dirs.sort(sortFn), ...files.sort(sortFn)]
})

// ── Navigation ──────────────────────────────────────────────────
function navigateTo(prefix) {
  currentPrefix.value = prefix
  searchQuery.value   = ''
  selected.value      = new Set()
  load()
}

function navigateUp() {
  if (!currentPrefix.value) return
  const t = currentPrefix.value.replace(/\/$/, '')
  navigateTo(t.includes('/') ? t.slice(0, t.lastIndexOf('/') + 1) : '')
}

// ── Data loading ────────────────────────────────────────────────
async function load() {
  loading.value     = true
  browseError.value = ''
  entries.value     = []
  nextPageToken.value = ''
  try {
    const result = await browseObjects(props.conn.provider, props.conn.bucket, props.conn.credentials, currentPrefix.value)
    entries.value       = result.entries ?? []
    nextPageToken.value = result.next_page_token ?? ''
  } catch (err) {
    browseError.value = err.message
  } finally {
    loading.value = false
    setupObserver()
  }
}

async function loadMore() {
  if (!nextPageToken.value || loadingMore.value) return
  loadingMore.value = true
  try {
    const result = await browseObjects(props.conn.provider, props.conn.bucket, props.conn.credentials, currentPrefix.value, nextPageToken.value)
    entries.value.push(...(result.entries ?? []))
    nextPageToken.value = result.next_page_token ?? ''
  } catch (err) {
    toast.error('Failed to load more: ' + err.message)
  } finally {
    loadingMore.value = false
  }
}

function refresh() { load() }

// ── Infinite scroll ─────────────────────────────────────────────
function setupObserver() {
  if (observer) { observer.disconnect(); observer = null }
  nextTick(() => {
    if (!sentinel.value) return
    observer = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting) loadMore()
    }, { root: bodyEl.value, threshold: 0.1 })
    observer.observe(sentinel.value)
  })
}

// ── Stats ───────────────────────────────────────────────────────
async function loadStats() {
  if (statsLoaded.value) return
  statsLoading.value = true
  statsError.value   = ''
  try {
    stats.value      = await getBucketStats(props.conn.provider, props.conn.bucket, props.conn.credentials)
    statsLoaded.value = true
  } catch (err) {
    statsError.value = 'Stats unavailable'
  } finally {
    statsLoading.value = false
  }
}

function toggleStats() {
  showStats.value = !showStats.value
  if (showStats.value && !statsLoaded.value) loadStats()
}

// ── Sort ────────────────────────────────────────────────────────
function cycleSort(key) {
  if (sortKey.value === key) sortDir.value = sortDir.value === 'asc' ? 'desc' : 'asc'
  else { sortKey.value = key; sortDir.value = 'asc' }
}

// ── Bulk operations ─────────────────────────────────────────────
async function bulkDelete() {
  const names = [...selected.value]
  const ok = await confirm.confirm(`Delete ${names.length} file${names.length > 1 ? 's' : ''}? This cannot be undone.`, 'Bulk Delete')
  if (!ok) return
  bulkWorking.value = true
  let failed = 0
  for (const name of names) {
    try {
      await deleteObject(props.conn.provider, props.conn.bucket, props.conn.credentials, name)
    } catch { failed++ }
  }
  selected.value = new Set()
  bulkWorking.value = false
  if (failed) toast.error(`${failed} file(s) could not be deleted.`)
  else toast.success(`${names.length} file${names.length > 1 ? 's' : ''} deleted.`)
  await load()
  if (statsLoaded.value) { statsLoaded.value = false; loadStats() }
}

async function bulkDownload() {
  const files = filteredEntries.value.filter(e => e.type === 'file' && selected.value.has(e.name))
  bulkWorking.value = true
  for (const entry of files) {
    try {
      const url = await getDownloadURL(props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name)
      const a = document.createElement('a')
      a.href = url; a.download = entry.display; a.target = '_blank'; a.rel = 'noopener'
      document.body.appendChild(a); a.click(); document.body.removeChild(a)
      await new Promise(r => setTimeout(r, 300)) // small delay to avoid popup block
    } catch (err) { toast.error('Download failed: ' + entry.display) }
  }
  bulkWorking.value = false
}

// ── Copy path ───────────────────────────────────────────────────
function copyPath(entry) {
  const path = props.conn.provider === 'gcp'
    ? `gs://${props.conn.bucket}/${entry.name}`
    : props.conn.provider === 'azure'
    ? `az://${props.conn.bucket}/${entry.name}`
    : props.conn.provider === 'gdrive'
    ? `gdrive://${props.conn.bucket}/${entry.name}`
    : `s3://${props.conn.bucket}/${entry.name}`
  navigator.clipboard?.writeText(path).then(
    () => toast.success('Path copied to clipboard'),
    () => toast.error('Clipboard not available'),
  )
}

// ── Download ────────────────────────────────────────────────────
async function download(entry) {
  try {
    const url = await getDownloadURL(props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name)
    const a = document.createElement('a')
    a.href = url; a.download = entry.display; a.target = '_blank'; a.rel = 'noopener'
    document.body.appendChild(a); a.click(); document.body.removeChild(a)
  } catch (err) {
    toast.error('Download failed: ' + err.message)
  }
}

// ── Delete single ───────────────────────────────────────────────
async function confirmDelete(entry) {
  const ok = await confirm.confirm(`Delete "${entry.display}"? This cannot be undone.`)
  if (!ok) return
  try {
    await deleteObject(props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name)
    if (previewEntry.value?.name === entry.name) closePreview()
    toast.success(`"${entry.display}" deleted.`)
    await load()
    if (statsLoaded.value) { statsLoaded.value = false; loadStats() }
  } catch (err) {
    toast.error('Delete failed: ' + err.message)
  }
}

// ── Upload ──────────────────────────────────────────────────────
async function handleUpload(files) {
  if (!files || files.length === 0) return
  uploading.value      = true
  uploadingCount.value = files.length
  try {
    await uploadObjects(props.conn.provider, props.conn.bucket, props.conn.credentials, currentPrefix.value, files)
    toast.success(`${files.length} file${files.length > 1 ? 's' : ''} uploaded.`)
    await load()
    if (statsLoaded.value) { statsLoaded.value = false; loadStats() }
  } catch (err) {
    toast.error('Upload failed: ' + err.message)
  } finally {
    uploading.value = false
  }
}

function onFileInput(e) { handleUpload(e.target.files); e.target.value = '' }

let dragCounter = 0
function onDragOver(e) { if (!e.dataTransfer?.types.includes('Files')) return; dragCounter++; isDragging.value = true }
function onDragLeave()  { if (--dragCounter <= 0) { dragCounter = 0; isDragging.value = false } }
function onDrop(e)      { dragCounter = 0; isDragging.value = false; handleUpload(e.dataTransfer?.files) }

// ── Create folder ────────────────────────────────────────────────
async function createFolder() {
  const name = newFolderName.value.trim()
  if (!name) return
  const prefix  = currentPrefix.value + name.replace(/\/+$/, '') + '/'
  const keepFile = new File([''], '.keep', { type: 'application/octet-stream' })
  showFolderModal.value = false
  newFolderName.value   = ''
  try {
    await uploadObjects(props.conn.provider, props.conn.bucket, props.conn.credentials, prefix, [keepFile])
    toast.success(`Folder "${name}" created.`)
    await load()
  } catch (err) {
    toast.error('Create folder failed: ' + err.message)
  }
}

watch(showFolderModal, open => { if (open) nextTick(() => folderInput.value?.focus()) })

// ── Rename / Move ────────────────────────────────────────────────
function openRename(entry) {
  renameEntry.value  = entry
  renameTarget.value = entry.display
  showRenameModal.value = true
  nextTick(() => renameInput.value?.focus())
}

async function doRename() {
  const target = renameTarget.value.trim()
  if (!target || renaming.value) return
  const destination = currentPrefix.value + target
  if (destination === renameEntry.value.name) { showRenameModal.value = false; return }
  renaming.value = true
  try {
    await copyObject(props.conn.provider, props.conn.bucket, props.conn.credentials, renameEntry.value.name, destination, true)
    if (previewEntry.value?.name === renameEntry.value.name) closePreview()
    toast.success(`Renamed to "${target}".`)
    showRenameModal.value = false
    await load()
  } catch (err) {
    toast.error('Rename failed: ' + err.message)
  } finally {
    renaming.value = false
  }
}

// ── Preview ─────────────────────────────────────────────────────
function isImage(entry) {
  const ct  = (entry?.content_type || '').toLowerCase()
  const ext = entry?.display.split('.').pop().toLowerCase()
  return ct.startsWith('image/') || ['jpg','jpeg','png','gif','webp','svg','ico','bmp'].includes(ext)
}
function isText(entry) {
  const ct  = (entry?.content_type || '').toLowerCase()
  const ext = entry?.display.split('.').pop().toLowerCase()
  return ct.startsWith('text/') || ct.includes('json') || ct.includes('xml') || ct.includes('javascript')
    || ['txt','md','json','yaml','yml','toml','csv','xml','html','js','ts','py','sh','log','conf','ini'].includes(ext)
}

async function openPreview(entry) {
  metaEntry.value = null
  previewEntry.value     = entry
  previewUrl.value       = ''
  previewContent.value   = ''
  previewLoadError.value = false
  previewLoading.value   = true
  try {
    const url = await getDownloadURL(props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name)
    previewUrl.value = url
    if (isText(entry)) {
      const res = await fetch(url)
      if (res.ok) previewContent.value = (await res.text()).slice(0, 50_000)
      else previewLoadError.value = true
    }
  } catch { previewLoadError.value = true }
  finally { previewLoading.value = false }
}

function closePreview() { previewEntry.value = null; previewUrl.value = ''; previewContent.value = '' }

// ── Metadata editor ─────────────────────────────────────────────
async function openMeta(entry) {
  previewEntry.value = null
  metaEntry.value    = entry
  metaData.value     = null
  metaError.value    = ''
  metaLoading.value  = true
  try {
    const data = await getObjectMetadata(props.conn.provider, props.conn.bucket, props.conn.credentials, entry.name)
    metaData.value = data
    metaEdit.value = { content_type: data.content_type || '', cache_control: data.cache_control || '' }
    metaRows.value = Object.entries(data.metadata || {}).map(([key, val]) => ({ key, val }))
  } catch (err) {
    metaError.value = 'Failed to load metadata: ' + err.message
  } finally {
    metaLoading.value = false
  }
}

function addMetaRow() { metaRows.value.push({ key: '', val: '' }) }

async function saveMeta() {
  metaSaving.value = true
  try {
    const metadata = {}
    for (const { key, val } of metaRows.value) {
      if (key.trim()) metadata[key.trim()] = val
    }
    await updateObjectMetadata(props.conn.provider, props.conn.bucket, props.conn.credentials, metaEntry.value.name, {
      content_type:  metaEdit.value.content_type,
      cache_control: metaEdit.value.cache_control,
      metadata,
    })
    toast.success('Metadata saved.')
    metaEntry.value = null
  } catch (err) {
    toast.error('Save failed: ' + err.message)
  } finally {
    metaSaving.value = false
  }
}

// ── Keyboard shortcuts ──────────────────────────────────────────
function onKeyDown(e) {
  const inInput = ['INPUT', 'TEXTAREA'].includes(document.activeElement?.tagName)
  if (e.key === '/' && !inInput)                                { e.preventDefault(); searchInput.value?.focus(); return }
  if ((e.key === 'r' || e.key === 'R') && !inInput && !e.metaKey && !e.ctrlKey) { refresh(); return }
  if (e.key === 'Escape') {
    if (metaEntry.value)    { metaEntry.value = null; return }
    if (previewEntry.value) { closePreview(); return }
    if (searchQuery.value)  { searchQuery.value = ''; return }
  }
  if (e.key === 'Backspace' && !inInput && currentPrefix.value) { e.preventDefault(); navigateUp() }
}

// ── Watchers / lifecycle ─────────────────────────────────────────
watch(() => props.conn, () => {
  currentPrefix.value = ''
  searchQuery.value   = ''
  selected.value      = new Set()
  stats.value         = null
  statsLoaded.value   = false
  previewEntry.value  = null
  metaEntry.value     = null
  load()
})

watch(sentinel, val => {
  if (observer) observer.disconnect()
  if (val && nextPageToken.value) {
    observer = new IntersectionObserver(es => { if (es[0].isIntersecting) loadMore() }, { root: bodyEl.value, threshold: 0.1 })
    observer.observe(val)
  }
})

onMounted(() => { load(); window.addEventListener('keydown', onKeyDown) })
onUnmounted(() => { window.removeEventListener('keydown', onKeyDown); observer?.disconnect() })

// ── Formatters ──────────────────────────────────────────────────
function formatSize(bytes) {
  if (!bytes && bytes !== 0) return '—'
  if (bytes === 0) return '0 B'
  const u = ['B','KB','MB','GB','TB']
  const i = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), u.length - 1)
  return (bytes / Math.pow(1024, i)).toFixed(i === 0 ? 0 : 1) + ' ' + u[i]
}
function formatDate(iso) {
  if (!iso) return '—'
  try { return new Date(iso).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' }) }
  catch { return '—' }
}
</script>
