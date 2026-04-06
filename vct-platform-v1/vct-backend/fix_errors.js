const fs = require('fs');
const path = require('path');

const dir = 'd:/VCT PLATFORM/vct-platform/backend/internal/httpapi';
const files = fs.readdirSync(dir).filter(f => f.endsWith('.go'));

for (const file of files) {
  const filePath = path.join(dir, file);
  let content = fs.readFileSync(filePath, 'utf8');
  let original = content;

  content = content.replace(/internalError\(/g, 'apiInternal(');
  content = content.replace(/methodNotAllowed\(/g, 'apiMethodNotAllowed(');
  content = content.replace(/notFoundError\(([^,]+),\s*([^)]+)\)/g, 'apiError($1, http.StatusNotFound, CodeNotFound, $2)');
  content = content.replace(/forbiddenErr\(/g, 'apiForbidden(');
  
  if (content !== original) {
    fs.writeFileSync(filePath, content, 'utf8');
    console.log('Fixed', file);
  }
}
