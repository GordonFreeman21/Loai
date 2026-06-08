<script>
  import { onMount, tick } from 'svelte';
  import Markdown from './lib/Markdown.svelte';
  import * as Events from '../wailsjs/runtime/runtime';
  import {
    GetModels,
    GetChats,
    GetMessages,
    CreateChat,
    SaveMessage,
    SendMessage,
    DeleteChat,
    PullModel
  } from '../wailsjs/go/backend/App';

  let models = [];
  let selectedModel = '';
  let chats = [];
  let currentChatId = null;
  let messages = [];
  let userInput = '';
  let isStreaming = false;
  let sidebarOpen = true;
  let newModelName = '';
  let showPullModal = false;
  let pullStatus = '';

  onMount(async () => {
    try {
      models = await GetModels();
      if (models.length > 0) selectedModel = models[0];
      await loadChats();
    } catch (e) {
      console.error(e);
    }

    Events.EventsOn("ollama_chunk", (chunk) => {
      if (messages.length > 0 && messages[messages.length - 1].role === 'assistant') {
        messages[messages.length - 1].content += chunk;
        messages = [...messages];
      } else {
        messages = [...messages, { role: 'assistant', content: chunk }];
      }
      scrollToBottom();
    });

    Events.EventsOn("ollama_done", async () => {
      isStreaming = false;
      const lastMessage = messages[messages.length - 1];
      await SaveMessage(currentChatId, lastMessage.role, lastMessage.content);
    });

    Events.EventsOn("pull_progress", (status) => {
      pullStatus = status;
    });

    Events.EventsOn("pull_done", async () => {
      pullStatus = 'Done!';
      models = await GetModels();
      setTimeout(() => { showPullModal = false; pullStatus = ''; }, 2000);
    });
  });

  async function loadChats() {
    chats = await GetChats();
  }

  async function selectChat(id) {
    currentChatId = id;
    messages = await GetMessages(id);
    scrollToBottom();
  }

  async function startNewChat() {
    currentChatId = null;
    messages = [];
  }

  async function handleSend() {
    if (!userInput.trim() || isStreaming) return;

    if (!currentChatId) {
      currentChatId = await CreateChat(userInput.substring(0, 30));
      await loadChats();
    }

    const userMsg = userInput;
    userInput = '';

    messages = [...messages, { role: 'user', content: userMsg }];
    await SaveMessage(currentChatId, 'user', userMsg);

    isStreaming = true;
    scrollToBottom();

    try {
      await SendMessage(currentChatId, selectedModel, messages);
    } catch (e) {
      isStreaming = false;
      messages = [...messages, { role: 'assistant', content: "Error: " + e }];
    }
  }

  async function handleDeleteChat(id, event) {
    event.stopPropagation();
    await DeleteChat(id);
    if (currentChatId === id) {
      currentChatId = null;
      messages = [];
    }
    await loadChats();
  }

  async function handlePullModel() {
    if (!newModelName.trim()) return;
    pullStatus = 'Starting...';
    await PullModel(newModelName);
    newModelName = '';
  }

  function scrollToBottom() {
    tick().then(() => {
      const container = document.getElementById('chat-container');
      if (container) {
        container.scrollTop = container.scrollHeight;
      }
    });
  }

  function handleKeydown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }
</script>

<div class="flex h-screen bg-[#343541] text-white overflow-hidden">
  <!-- Sidebar -->
  {#if sidebarOpen}
    <div class="w-[260px] bg-[#202123] flex flex-col h-full">
      <button
        on:click={startNewChat}
        class="m-2 p-3 border border-white/20 rounded hover:bg-white/5 transition-colors text-left flex items-center gap-3"
      >
        <span class="text-xl">+</span> New Chat
      </button>

      <div class="flex-1 overflow-y-auto px-2 space-y-2 mt-2">
        {#each chats as chat}
          <div
            on:click={() => selectChat(chat.id)}
            class="p-3 rounded hover:bg-[#2A2B32] cursor-pointer group flex justify-between items-center {currentChatId === chat.id ? 'bg-[#2A2B32]' : ''}"
          >
            <span class="truncate text-sm">{chat.title}</span>
            <button
              on:click={(e) => handleDeleteChat(chat.id, e)}
              class="opacity-0 group-hover:opacity-100 hover:text-red-500 p-1"
            >
              🗑
            </button>
          </div>
        {/each}
      </div>

      <div class="p-2 border-t border-white/10">
        <button
          on:click={() => showPullModal = true}
          class="w-full p-3 rounded hover:bg-[#2A2B32] text-left text-sm"
        >
          ⬇️ Pull New Model
        </button>
      </div>
    </div>
  {/if}

  <!-- Main Content -->
  <div class="flex-1 flex flex-col relative">
    <!-- Header -->
    <header class="h-14 border-b border-white/10 flex items-center justify-between px-4">
      <button on:click={() => sidebarOpen = !sidebarOpen} class="p-2 hover:bg-white/5 rounded">
        ☰
      </button>

      <div class="flex items-center gap-2">
        <select
          bind:value={selectedModel}
          class="bg-[#343541] border border-white/20 rounded px-3 py-1 text-sm outline-none focus:border-white/40"
        >
          {#each models as model}
            <option value={model}>{model}</option>
          {/each}
        </select>
        {#if models.length === 0}
          <span class="text-xs text-red-400">Ollama not detected</span>
        {/if}
      </div>
      <div class="w-10"></div>
    </header>

    <!-- Chat Area -->
    <main id="chat-container" class="flex-1 overflow-y-auto pb-32">
      {#if messages.length === 0}
        <div class="h-full flex items-center justify-center text-white/40">
          <div class="text-center">
            <h1 class="text-4xl font-bold mb-4">Ollama Chat</h1>
            <p>Select a model and start chatting locally.</p>
          </div>
        </div>
      {:else}
        {#each messages as msg}
          <div class="py-8 border-b border-black/10 {msg.role === 'assistant' ? 'bg-[#444654]' : ''}">
            <div class="max-w-3xl mx-auto px-4 flex gap-6">
              <div class="w-8 h-8 rounded flex items-center justify-center font-bold text-sm {msg.role === 'assistant' ? 'bg-[#10a37f]' : 'bg-[#5436da]'}">
                {msg.role === 'assistant' ? 'AI' : 'U'}
              </div>
              <div class="flex-1 min-w-0">
                <Markdown content={msg.content} />
              </div>
            </div>
          </div>
        {/each}
      {/if}
    </main>

    <!-- Input Area -->
    <div class="absolute bottom-0 left-0 right-0 p-4 bg-gradient-to-t from-[#343541] via-[#343541] to-transparent">
      <div class="max-w-3xl mx-auto relative">
        <textarea
          bind:value={userInput}
          on:keydown={handleKeydown}
          placeholder="Send a message..."
          rows="1"
          class="w-full bg-[#40414f] text-white border border-black/10 rounded-xl px-4 py-4 pr-12 shadow-2xl outline-none focus:ring-1 focus:ring-white/20 resize-none max-h-60"
          style="overflow-y: hidden;"
          on:input={(e) => {
            e.target.style.height = 'auto';
            e.target.style.height = e.target.scrollHeight + 'px';
          }}
        ></textarea>
        <button
          on:click={handleSend}
          disabled={!userInput.trim() || isStreaming}
          class="absolute right-2 bottom-3 p-2 rounded-md transition-colors {userInput.trim() && !isStreaming ? 'bg-[#10a37f] text-white' : 'text-white/20'}"
        >
          ➤
        </button>
      </div>
      <p class="text-center text-[10px] text-white/40 mt-2">
        Free Research Preview. LLMs may produce inaccurate information about people, places, or facts.
      </p>
    </div>
  </div>
</div>

<!-- Pull Model Modal -->
{#if showPullModal}
  <div class="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
    <div class="bg-[#202123] p-6 rounded-lg w-96 border border-white/10 shadow-2xl">
      <h2 class="text-xl font-bold mb-4">Pull New Model</h2>
      <input
        bind:value={newModelName}
        placeholder="e.g. llama3, mistral"
        class="w-full bg-[#343541] border border-white/20 rounded px-4 py-2 mb-4 outline-none focus:border-white/40"
      />
      {#if pullStatus}
        <div class="text-xs text-white/60 mb-4 font-mono whitespace-pre-wrap max-h-20 overflow-y-auto">
          {pullStatus}
        </div>
      {/if}
      <div class="flex justify-end gap-3">
        <button
          on:click={() => { showPullModal = false; pullStatus = ''; }}
          class="px-4 py-2 hover:bg-white/5 rounded transition-colors"
        >
          Cancel
        </button>
        <button
          on:click={handlePullModel}
          class="px-4 py-2 bg-[#10a37f] hover:bg-[#1a7f64] rounded transition-colors"
          disabled={!newModelName.trim() || (pullStatus && pullStatus !== 'Done!')}
        >
          Pull
        </button>
      </div>
    </div>
  </div>
{/if}

<style>
  :global(body) {
    overflow: hidden;
  }
</style>
