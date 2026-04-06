const fs = require('fs');
const path = require('path');

const dir = 'd:/VCT PLATFORM/vct-platform/backend/internal/httpapi';
const files = fs.readdirSync(dir).filter(f => f.endsWith('.go'));

for (const file of files) {
  const filePath = path.join(dir, file);
  let content = fs.readFileSync(filePath, 'utf8');
  let original = content;

  if (file !== 'helpers.go') {
    // Only replace whole words to avoid replacing inside apiUnauthorized or apiNotFound!
    
    // unauthorized(w, msg) -> apiUnauthorized(w, msg)
    content = content.replace(/\bunauthorized\(([^,]+),\s*([^)]+)\)/g, 'apiUnauthorized($1, $2)');
    
    // notFound(w) -> apiError(w, http.StatusNotFound, CodeNotFound, "Không tìm thấy tài nguyên")
    content = content.replace(/\bnotFound\(([^),]+)\)/g, 'apiError($1, http.StatusNotFound, CodeNotFound, "Không tìm thấy tài nguyên")');
    
    // conflict(w, msg) -> apiConflict(w, msg)
    content = content.replace(/\bconflict\(([^,]+),\s*([^)]+)\)/g, 'apiConflict($1, $2)');
  }

  if (content !== original) {
    fs.writeFileSync(filePath, content, 'utf8');
    console.log('Fixed', file);
  }
}
