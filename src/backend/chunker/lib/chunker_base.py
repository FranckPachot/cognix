import os
import logging
from langchain_text_splitters import RecursiveCharacterTextSplitter
from gen_types.chunking_data_pb2 import ChunkingData, FileType
from typing import List, Tuple
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

chunk_size = int(os.getenv('CHUNK_SIZE', 500))
chunk_overlap = int(os.getenv('CHUNK_OVEVRLAP', 3))

CHUNK_OVEVRLAP=20


class BaseChunker:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)

    def chunk(self, data: ChunkingData):
        raise NotImplementedError("Chunk method needs to be implemented by subclasses")
    
    def split_data(self, content: str, url: str) -> List[Tuple[str, str]]:
        # This method should split the content into chunks and return a list of tuples (chunk, url)
        # For demonstration, let's split content by lines
        logging.warning("split_data shall implement various chunk thechniques and compare them")
        

        # Initialize the text splitter with custom parameters
        custom_text_splitter = RecursiveCharacterTextSplitter(
            # Set custom chunk size
            chunk_size = chunk_size,
            chunk_overlap  = chunk_overlap,
            # Use length of the text as the size measure
            length_function = len,
            # Use only "\n\n" as the separator
            separators = ['\n']
        )

        # Create the chunks
        texts = custom_text_splitter.create_documents([content])

        self.logger.info(f"otiginal text {content}")
        self.logger.info("====   Sample chunks from 'Custom Parameters':   ====\n\n")
        for i, chunk in enumerate(texts):
            self.logger.info(f"### Chunk {i+1}: \n{chunk.page_content}\n")

        # TODO: add to the return list only if the chunk is longer than x chars
        return [(chunk.page_content, url) for chunk in texts if chunk]


