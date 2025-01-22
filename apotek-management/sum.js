const fs = require("fs");
const path = require("path");

const directoryPath = path.join(__dirname, "models");

const outputFileName = "./model.txt";

function readFilesRecursively(
  dir,
  ignored = ["node_modules", "package-lock.json"]
) {
  let filesContent = "";

  const files = fs.readdirSync(dir);

  for (const file of files) {
    const fullPath = path.join(dir, file);
    const stat = fs.statSync(fullPath);

    if (ignored.includes(file) || file.startsWith(".")) continue;

    if (stat.isDirectory()) {
      filesContent += readFilesRecursively(fullPath, ignored);
    } else if (path.extname(file) === ".go") {
      const content = fs.readFileSync(fullPath, "utf-8");
      filesContent += `// File: ${fullPath}\n`;
      filesContent += content + "\n\n";
    }
  }

  return filesContent;
}

try {
  console.log("Menggabungkan file .js...");
  const mergedContent = readFilesRecursively(directoryPath);

  fs.writeFileSync(outputFileName, mergedContent, "utf-8");
  console.log(`Semua file .js berhasil digabungkan ke dalam ${outputFileName}`);
} catch (error) {
  console.error("Terjadi kesalahan:", error);
}
