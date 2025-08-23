<!--suppress JSUnresolvedReference -->
<script>
    let library = {artists: []};
    let selectedArtist = null;
    let selectedAlbum = null;
    let nowPlaying = null;
    let isPlaying = false;
    let audioEl;
    let volume = 0.8;
    let duration = 0;
    let current = 0;
    let seeking = false;

    async function load() {
        const res = await fetch('/api/library');
        library = await res.json();
        if(library.artists && library.artists.length) {
            selectArtist(library.artists[0]);
        }
    }

    function selectArtist(artist) {
        selectedArtist = artist;
        selectedAlbum = artist && artist.albums && artist.albums[0] || null;
    }

    function selectAlbum(album) {
        selectedAlbum = album;
    }

    function playTrack(track) {
        nowPlaying = track;
        const url = '/stream/' + track.id;
        if(!audioEl) return;
        if(audioEl.src !== location.origin + url) {
            audioEl.src = url;
        }
        audioEl.play();
        isPlaying = true;
    }

    function togglePlay() {
        if(!audioEl) return;
        if(audioEl.paused) {
            audioEl.play();
            isPlaying = true;
        } else {
            audioEl.pause();
            isPlaying = false;
        }
    }

    function startSeek() {
        if(!audioEl) return;
        seeking = true;
    }

    function moveSeek(event) {
        if(!seeking) return;
        current = +event.target.value;
    }

    function endSeek(event) {
        if(!audioEl || !seeking) return;
        const val = +event.target.value;
        audioEl.currentTime = val;
        current = val;
        seeking = false;
    }

    async function rescan() {
        await fetch('/api/rescan', {method: 'POST'});
        await load();
    }

    function onAudioPlay() {
        isPlaying = true;
    }

    function onAudioPause() {
        isPlaying = false;
    }

    function onVolume(event) {
        volume = +event.target.value;
        if(audioEl) audioEl.volume = volume;
    }

    function onLoadedMetadata() {
        duration = isFinite(audioEl.duration) ? audioEl.duration : 0;
    }

    function onTimeUpdate() {
        if(audioEl && !seeking) current = audioEl.currentTime;
    }

    if(audioEl) {
        audioEl.volume = volume;
    }

    load();
</script>

<svelte:window on:keydown={(e) => { if(e.code==='Space'){ e.preventDefault(); togglePlay(); } }}/>

<div class="layout">
    <div class="sidebar">
        <div class="brand" style="padding: 10px 8px 6px;">
            <span class="dot"></span>
            <span>Sonis</span>
        </div>
        <h2>Artists</h2>
        <div class="artist-list">
            {#each library.artists as artist}
                <button
                        type="button"
                        class="artist {selectedArtist && selectedArtist.name===artist.name ? 'active' : ''}"
                        aria-pressed={selectedArtist && selectedArtist.name===artist.name}
                        on:click={() => selectArtist(artist)}
                        title={artist.name}>
                    {artist.name}
                </button>
            {/each}
        </div>
    </div>

    <div class="content">
        <div class="header">
            <div><strong>{selectedArtist ? selectedArtist.name : 'Library'}</strong></div>
            <button class="rescan" type="button" on:click={rescan} title="Rescan library">
                <i class="fa-solid fa-rotate" aria-hidden="true"></i>
                <span class="visually-hidden">Rescan</span>
            </button>
        </div>

        <div class="main">
            <div class="albums">
                <div class="section-title"><i class="fa-regular fa-rectangle-list" aria-hidden="true"></i> Albums</div>
                <div class="album-grid">
                    {#if selectedArtist}
                        {#each selectedArtist.albums as album}
                            <button class="album" type="button" on:click={() => selectAlbum(album)} title={album.name}>
                                <div class="album-name">{album.name}</div>
                                <div class="album-artist">{album.artist}</div>
                            </button>
                        {/each}
                    {:else}
                        <div>No artists found</div>
                    {/if}
                </div>
            </div>

            <div class="tracks">
                <div class="section-title">
                    <i class="fa-solid fa-music" aria-hidden="true"></i>
                    Tracks
                    {#if selectedAlbum}— {selectedAlbum.name}{/if}
                </div>
                {#if selectedAlbum}
                    <ul class="track-list">
                        {#each selectedAlbum.tracks as t, i}
                            <li class="track {nowPlaying && nowPlaying.id===t.id ? 'playing' : ''}">
                                <button
                                        type="button"
                                        class="track-btn"
                                        on:click={() => playTrack(t)}
                                        title={`Play ${t.name}`}>
                                    <div class="track-idx">
                                        {#if nowPlaying && nowPlaying.id === t.id}
                                            <i class="fa-solid fa-volume-high" aria-hidden="true"></i>
                                        {:else}
                                            {t.track || i + 1}
                                        {/if}
                                    </div>
                                    <div class="track-name">{t.name}</div>
                                    <div class="track-album">{t.album}</div>
                                </button>
                            </li>
                        {/each}
                    </ul>
                {:else}
                    <div style="color: var(--muted);">Select an album to see tracks</div>
                {/if}
            </div>
        </div>
    </div>
</div>

<div class="player">
    <div class="now-playing">{nowPlaying ? `${nowPlaying.artist} — ${nowPlaying.name}` : 'Nothing playing'}</div>
    <div class="controls">
        <button class="play-btn" type="button" on:click={togglePlay} title={isPlaying ? 'Pause' : 'Play'}>
            {#if isPlaying}
                <i class="fa-solid fa-pause" aria-hidden="true"></i>
            {:else}
                <i class="fa-solid fa-play" aria-hidden="true"></i>
            {/if}
            <span class="visually-hidden">{isPlaying ? 'Pause' : 'Play'}</span>
        </button>
        <input
                class="seek"
                type="range"
                min="0"
                max={duration}
                step="0.1"
                value={current}
                on:mousedown={startSeek}
                on:touchstart={startSeek}
                on:input={moveSeek}
                on:mouseup={endSeek}
                on:touchend={endSeek}
                on:change={endSeek}
        />
    </div>
    <div class="volume">
        <i class="fa-solid fa-volume-high" aria-hidden="true" style="color: var(--muted); margin-right: 8px;"></i>
        <input class="vol-range" type="range" min="0" max="1" step="0.01" value={volume} on:input={onVolume}
               aria-label="Volume"/>
    </div>
</div>

<audio bind:this={audioEl}
       on:play={onAudioPlay}
       on:pause={onAudioPause}
       on:loadedmetadata={onLoadedMetadata}
       on:timeupdate={onTimeUpdate}
       preload="metadata"></audio>
