#!/usr/bin/env bash
set -euo pipefail

# 1) Folder for your toy corpus
mkdir -p corpus

# 2) Pick 10‒20 article titles you’d like to index
titles=(
  Earth
  Sun
  Moon
  Cat
  Poker
  "Python_(programming_language)"
  Serbia
  Rocket
  Computer
  "Artificial_intelligence"
)

# 3) Download each article’s fully rendered HTML
for title in "${titles[@]}"; do
  # Sanitise the filename (slashes/spaces → underscores)
  clean=$(echo "$title" | tr ' /()' '____')
  
  curl -s -L \
    -H "User-Agent: search-mvp/0.1 (you@example.com)" \
    "https://api.wikimedia.org/core/v1/wikipedia/en/page/${title}/html" \
    -o "corpus/${clean}.html"

  echo "✓  Saved ${clean}.html"
done
