from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class SearchRequest(_message.Message):
    __slots__ = ("query", "page", "count")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    PAGE_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    query: str
    page: int
    count: int
    def __init__(self, query: _Optional[str] = ..., page: _Optional[int] = ..., count: _Optional[int] = ...) -> None: ...

class SearchResponse(_message.Message):
    __slots__ = ("results", "total")
    RESULTS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    results: _containers.RepeatedCompositeFieldContainer[SearchResult]
    total: int
    def __init__(self, results: _Optional[_Iterable[_Union[SearchResult, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class SearchResult(_message.Message):
    __slots__ = ("Doc", "Score", "TermCount")
    DOC_FIELD_NUMBER: _ClassVar[int]
    SCORE_FIELD_NUMBER: _ClassVar[int]
    TERMCOUNT_FIELD_NUMBER: _ClassVar[int]
    Doc: DocMetadata
    Score: float
    TermCount: int
    def __init__(self, Doc: _Optional[_Union[DocMetadata, _Mapping]] = ..., Score: _Optional[float] = ..., TermCount: _Optional[int] = ...) -> None: ...

class DocMetadata(_message.Message):
    __slots__ = ("url", "depth", "title", "hash", "images", "first_paragraph")
    URL_FIELD_NUMBER: _ClassVar[int]
    DEPTH_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    HASH_FIELD_NUMBER: _ClassVar[int]
    IMAGES_FIELD_NUMBER: _ClassVar[int]
    FIRST_PARAGRAPH_FIELD_NUMBER: _ClassVar[int]
    url: str
    depth: int
    title: str
    hash: str
    images: _containers.RepeatedScalarFieldContainer[str]
    first_paragraph: str
    def __init__(self, url: _Optional[str] = ..., depth: _Optional[int] = ..., title: _Optional[str] = ..., hash: _Optional[str] = ..., images: _Optional[_Iterable[str]] = ..., first_paragraph: _Optional[str] = ...) -> None: ...
