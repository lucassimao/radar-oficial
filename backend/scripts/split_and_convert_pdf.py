import pymupdf4llm
import pymupdf as PyMuPDF
import tempfile
import sys
import os
from pathlib import Path

# Split a PDF file into pages and convert each page to Markdown
# Each final Markdown file will be named following the pattern:
#
# page_1.md
# page_2.md
# page_X.md
#
# Where X is the total of pages in the original PDF file
#
# Usage:
#
# python split_and_convert_pdf.py << PATH TO PDF FILE >> << PATH TO OUTPUT DIR >>

if len(sys.argv) == 1:
    sys.exit("Use: split_and_convert_pdf.py << PATH TO PDF FILE >> << PATH TO OUTPUT DIR >>")

# Define paths
input_pdf_path =  sys.argv[1]

if not os.path.isfile(input_pdf_path):
    sys.exit("Input file does not exist")
    
# Throws an error if anything goes wrong
output_dir = Path(sys.argv[2])
output_dir.mkdir(exist_ok=True)

# Open the PDF to get page count
doc = PyMuPDF.open(input_pdf_path)
total_pages = len(doc)

print(f"Processing PDF with {total_pages} pages...")

# Process each page separately
for page_num in range(total_pages):
    # Create a temporary file that deletes on close
    with tempfile.NamedTemporaryFile(delete=True) as temp_pdf_page:
        
        new_doc = PyMuPDF.open()
        new_doc.insert_pdf(doc, from_page=page_num, to_page=page_num)
        new_doc.save(temp_pdf_page.name)
        new_doc.close()

        
        print(f"Processing page {page_num + 1}/{total_pages}...")
        # Convert the page to markdown
        try:
            md_text = pymupdf4llm.to_markdown(str(temp_pdf_page.name))
            
            # Save the full page as markdown
            md_path = output_dir / f"page_{page_num + 1:03}.md"
            with md_path.open("w", encoding="utf-8") as f:
                f.write(md_text)
        except Exception as e:
            print(f"Error processing page {page_num + 1}: {e}")

doc.close()
print(f"All pages processed. Results saved to {output_dir}")


# Function to extract a single page as a separate PDF
def extract_page_to_pdf(input_path, page_num, output_path):
    doc = PyMuPDF.open(input_path)
    new_doc = PyMuPDF.open()
    new_doc.insert_pdf(doc, from_page=page_num, to_page=page_num)
    new_doc.save(output_path)
    new_doc.close()
    doc.close()
