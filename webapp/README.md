# Webapp (Frontend)

A modern React-based search interface built with Vite, TypeScript, and Tailwind CSS that provides a Google-like search experience with autocomplete, spell suggestions, and responsive design.

## Overview

The Webapp is the user-facing frontend of the search engine, providing an intuitive and responsive search interface. Built with modern React patterns and optimized for performance, it offers features like real-time autocomplete, "I'm Feeling Lucky" functionality, pagination, and image search results.

## Architecture

### Core Components

- **React 18**: Modern React with hooks and concurrent features
- **TypeScript**: Type-safe development with full type coverage
- **Vite**: Fast build tool with hot module replacement
- **Tailwind CSS**: Utility-first CSS framework for rapid styling
- **React Router**: Client-side routing for search results
- **API Integration**: RESTful communication with Query API service

## Key Features

### 1. Google-like Search Interface
- **Clean homepage**: Minimalist design with centered search box
- **Logo branding**: Custom search engine branding
- **Responsive design**: Mobile-first responsive layout
- **Keyboard shortcuts**: Enter to search, escape to clear

### 2. Advanced Search Bar
- **Real-time autocomplete**: Live suggestions as user types
- **Debounced input**: Optimized API calls with input debouncing
- **Keyboard navigation**: Arrow keys and enter for suggestion selection
- **Visual feedback**: Hover and selection states for suggestions

### 3. Search Results Display
- **Ranked results**: BM25-scored results with relevance indicators
- **Rich snippets**: First paragraph excerpts for result preview
- **Metadata display**: Title, URL, and content information
- **Query time**: Display of search processing time

### 4. Intelligent Features
- **"I'm Feeling Lucky"**: Direct navigation to top result
- **Spell correction**: "Did you mean?" suggestions for misspelled queries
- **Pagination**: Efficient result pagination with page navigation
- **Image search**: Dedicated image results view

## Performance Characteristics

### Load Times
- **Initial load**: <2 seconds on 3G connection
- **Route transitions**: <100ms for client-side navigation
- **Search results**: <500ms including API call
- **Autocomplete**: <50ms for suggestion display

### Bundle Size
- **Main bundle**: ~200KB gzipped
- **Vendor bundle**: ~150KB gzipped (React, React Router)
- **CSS bundle**: ~50KB gzipped (Tailwind CSS)
- **Total size**: ~400KB gzipped for initial load